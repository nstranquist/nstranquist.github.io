.PHONY: render check typecheck test verify serve

render:
	go run ./tools/render --root .

check:
	go run ./tools/render --root . --check

typecheck:
	go vet ./tools/render

test:
	go test ./tools/render

verify: typecheck test render check

serve: render
	python3 -m http.server 8766 --bind 127.0.0.1
