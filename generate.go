package generate

//go:generate tools/bin/mockgen -source=internal/processor/port/port.go -destination=internal/processor/port/port_mock.go -package=port

//go:generate tools/bin/mockgen -source=internal/processor/port/msgsender.go -destination=internal/processor/port/msgsender_mock.go -package=port

//go:generate tools/bin/mockgen -source=internal/processor/port/rankings.go -destination=internal/processor/port/rankings_mock.go -package=port
