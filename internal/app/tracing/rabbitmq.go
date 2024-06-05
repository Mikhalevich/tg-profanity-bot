package tracing

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type headersCarrier map[string]interface{}

func (h headersCarrier) Get(key string) string {
	v, ok := h[key]
	if !ok {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}

func (h headersCarrier) Set(key string, value string) {
	h[key] = value
}

func (h headersCarrier) Keys() []string {
	keys := make([]string, 0, len(h))

	for k := range h {
		keys = append(keys, k)
	}

	return keys
}

// injectHeaders injects the tracing from the context into the header map.
func injectContextHeaders(ctx context.Context) map[string]interface{} {
	h := make(headersCarrier)
	otel.GetTextMapPropagator().Inject(ctx, h)

	return h
}

// extractHeaders extracts the tracing from the header and puts it into the context.
func ExtractContextFromHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, headersCarrier(headers))
}

type channelWrapper struct {
	*amqp.Channel
}

func WrapChannel(channel *amqp.Channel) *channelWrapper {
	return &channelWrapper{
		Channel: channel,
	}
}

func (cw *channelWrapper) PublishWithContext(
	ctx context.Context,
	exchange string,
	key string,
	mandatory bool,
	immediate bool,
	msg amqp.Publishing,
) error {
	traceHeaders := injectContextHeaders(ctx)
	if len(traceHeaders) > 0 {
		if msg.Headers == nil {
			msg.Headers = make(amqp.Table, len(traceHeaders))
		}

		for k, v := range traceHeaders {
			msg.Headers[k] = v
		}
	}

	if err := cw.Channel.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg); err != nil {
		return fmt.Errorf("channel wrapper PublishWithContext: %w", err)
	}

	return nil
}
