.PHONY: render check test verify serve

render:
	go run ./tools/render --root .

check:
	go run ./tools/render --root . --check

test:
	go test ./tools/render

verify: test render check

serve: render
	python3 -m http.server 8766 --bind 127.0.0.1
