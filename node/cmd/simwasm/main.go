//go:build js && wasm

// Command simwasm compila o núcleo do protocolo (swap.NextAction + simulador) para WebAssembly.
// O modo offline do simulador web usa exatamente a mesma lógica de decisão dos nós Janus.
//
//	GOOS=js GOARCH=wasm go build -o apps/web/public/janus.wasm ./cmd/simwasm
//
// API exposta em globalThis:
//
//	janusScenarios(): string            // JSON do catálogo de cenários
//	janusRun(configJSON: string): string // JSON de sim.Result (ou {"error": "..."})
package main

import (
	"encoding/json"
	"syscall/js"

	"janus/node/internal/swap/sim"
)

func jsonOrError(v any, err error) string {
	if err != nil {
		b, _ := json.Marshal(map[string]string{"error": err.Error()})
		return string(b)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return `{"error":"marshal"}`
	}
	return string(b)
}

func main() {
	js.Global().Set("janusScenarios", js.FuncOf(func(js.Value, []js.Value) any {
		return jsonOrError(sim.Scenarios, nil)
	}))
	js.Global().Set("janusRun", js.FuncOf(func(_ js.Value, args []js.Value) any {
		if len(args) != 1 {
			return jsonOrError(nil, errArgs)
		}
		var cfg sim.Config
		if err := json.Unmarshal([]byte(args[0].String()), &cfg); err != nil {
			return jsonOrError(nil, err)
		}
		r, err := sim.Run(cfg)
		return jsonOrError(r, err)
	}))
	js.Global().Set("janusReady", true)
	select {} // mantém o runtime vivo
}

type argsError struct{}

func (argsError) Error() string { return "janusRun(configJSON) expects one argument" }

var errArgs = argsError{}
