---
title: "Using Language Servers with CodeMirror 6"
date: 2021-03-11T15:36:59+06:00
draft: true
---

CodeMirror 6, a rewrite of the [CodeMirror](https://codemirror.net/) editor, brings quite a huge number of improvements. [Toph](https://toph.co/) has been using CodeMirror as the integrated code editor ever since the feature was introduced.

As CodeMirror 6 reached a stable interface with the promise of better touchscreen support, it was time for an upgrade! During which I wanted to introduce language server support.

The goal was to provide code completion, diagnostics, and hover tooltips. And, CodeMirror 6 makes it really easy to do all three.

All of these have been packaged into a small library and made available on NPM:

[ [Language Server Plugin for CodeMirror 6](https://www.npmjs.com/package/codemirror-languageserver) , [GitHub](https://github.com/FurqanSoftware/codemirror-languageserver) ]

## Code Completion

The `@codemirror/autocomplete` package provides an [`autocompletion`](https://codemirror.net/6/docs/ref/#autocomplete.autocompletion) extension. 

``` ts {linenos=false}
autocompletion(config⁠?: Object = {}) → Extension
```

CodeMirror already provides a UX for code completion. Given the context of where code completion is activated, all you need to do is provide the options for completion.

This can be done by providing one or more completion sources through the `override` property of the `config` object.

``` js
autocompletion({
	override: [
		(context) => {
			let triggered = false; // Set to true if the cursor is positioned right after a trigger character (from server capabilities). 
			let triggerCharacter = void 0; // If triggered, set this to the trigger character.
			if (triggered || context.explicit) {
				return requestCompletion(context, {
					kind: triggered ? 2 : 1
					character: triggerCharacter
				});
			}
		}
	]
})


class Plugin {
	// ...

	requestCompletion(context, trigger) {
		this.sendChange(/* ... */);

		return return this.client.request({
			method: 'textDocument/completion',
			params: { /* ... */ }
		}, timeout).then((result) => {
			let completions; // Transform result.items to CodeMirror's Completion objects.
			return completions;
		});
	}
	
	// ...
}
```

The `requestCompletion` function should return a Promise of [`CompletionResult`](https://codemirror.net/6/docs/ref/#autocomplete.CompletionResult).

## Diagnostics

Similar to the `autocompletion` extension, CodeMirror makes showing diagnostics easy by providing the [`linter`](https://codemirror.net/6/docs/ref/#lint.linter) extension through the `@codemirror/lint` package.

``` ts
linter(
    source: fn(view: EditorView) → readonly Diagnostic[] | Promise<readonly Diagnostic[]>
) → Extension
```

The extension calls the `source` function when changes are made to the code but only after the editor has become idle. 

However, due to how language servers publish diagnostics, a bit of a workaround is necessary. You see, when you send code changes to a language server, you do not get diagnostics immediately. The language server publishes a notification when diagnostics are ready.

When the linter requests for diagnostics, return a promise:

``` js
linter((view) => {
	return plugin.requestDiagnostics(view);
})

class Plugin {
	// ...

	requestDiagnostics(view) {
		this.sendChange(/* ... */);

		return new Promise((fulfill, reject) => {
			this.promises.push({
				type: 'diagnostics',
				fulfill: fulfill,
				reject: reject
			});
		});
	}

	// ...
}
```

And, fulfill it once a publishDiagnostics notification is received from the language server.

``` js
processDiagnostics({params}) {
	let diagnostics; // Transform params.diagnostics to CodeMirror's Diagnostic objects.
	this.promises = this.promises.filter((p) => {
		if (p.type === 'diagnostics') {
			p.fulfill(diagnostics);
			return false;
		} else {
			return true;
		}
	});
}
```

## Hover Tooltips

A huge thanks to Marijn for taking care of a [feature request](https://discuss.codemirror.net/t/return-promise-tooltip-from-hovertooltips-source-function/2967) in \~4 days. This made implementing hover tooltips as easy as code completion.

The [`hoverTooltip`](https://codemirror.net/6/docs/ref/#tooltip.hoverTooltip) extension from the `@codemirror/tooltip` package allows you to return a Promise of [`Tooltip`](https://codemirror.net/6/docs/ref/#tooltip.Tooltip).

``` ts
hoverTooltip(
    source: fn(
        view: EditorView,
        pos: number,
        side: -1 | 1
    ) → Tooltip | Promise<Tooltip | null> | null,
    options⁠?: {hideOnChange⁠?: boolean} = {}
) → Extension
```

And, the drill here is pretty simple: Once the source function is called by CodeMirror (which happens when a the mouse cursor is hovering a bit of code), we send a request to the language server and wait for it to respond with any relevant documentations. 

``` js
hoverTooltip((view, pos, side) => {
	return plugin.requestHoverTooltip(view, offsetToPos(view.state.doc, pos));
})

class Plugin {
	// ...

	requestHoverTooltip(view, pos) {
		this.sendChange(/* ... */);

		return new Promise((fulfill, reject) => {
			this.client.request({
				method: 'textDocument/hover',
				params: { /* ... */ }
			}, timeout).then((result) => {
				if (!result) {
					return null;
				}

				let pos = posToOffset(view.state.doc, range.start);
				let end = posToOffset(view.state.doc, range.end);
				let dom = document.createElement('div');
				dom.textContent = formatContents(result.contents);
				fulfill({
					pos: pos,
					end: end,
					create: function(view) {
						return {
							dom
						};
					},
					above: true
				});
			}).catch(function(reason) {
				reject(reason);
			});
		});
	}


	// ...
}
```

## Language Server over WebSocket

If you want to quickly try serving Language Server features over WebSocket, pick one of these:

- [jsonrpc-ws-proxy](https://www.npmjs.com/package/jsonrpc-ws-proxy)
- [lsp-ws-proxy](https://github.com/qualified/lsp-ws-proxy)

If you want the details, continue reading.

Language Servers speak JSON-RPC 2.0 over standard IO. To invoke a method or send a notification to a language server, you can write to the process's standard input.

Start a language server (in this example it is the language server for Go) in your terminal:

``` sh {linenos=false}
~ » gopls
```

And enter the following as input:

``` text
Content-Length: 61

{"jsonrpc":"2.0","id":"0","method":"initialize","params":{}}
```

A JSON-RPC request contains some headers (at least Content-Length) followed by an empty line, followed by the payload.

The language server process responds by writing to standard output:

``` text
Content-Length: 2168

{"jsonrpc":"2.0","result":{...},"id":"0"}
``` 

Notifications are similar, except that you do not expect any response for them.

You can now write a small daemon program that listens for WebSocket connections, and spins up a language server when a connection is established.

The program should then read incoming messages from the WebSocket. Since the message will only contain the payload, and not the headers, the program should first write the headers to the standard input of the language server, followed by the payload.

The program should also read from the language server's standard output, validate the headers before discarding them, and send the payload back through the WebSocket.

A trivial example of the above would look like this:

``` golang
// Adapted from https://github.com/gorilla/websocket/tree/master/examples/command.
// Error handling has been omitted for brevity.

func serveWs(w http.ResponseWriter, r *http.Request) {
	stack := r.URL.Query().Get("stack")

	ws, _ := upgrader.Upgrade(w, r, nil)
	defer ws.Close()

	// Start Language Server inside Docker using locally available tagged images.
	cmd := exec.Command("docker", "run", "-i", "lsp-"+stack)
	inw, _ := cmd.StdinPipe()
	outr, _ := cmd.StdoutPipe()
	cmd.Start()

	done := make(chan struct{})
	go pumpStdout(ws, outr, done) // Read from stdout, write to WebSocket.
	go ping(ws, done)

	pumpStdin(ws, inw) // Read from WebSocket, write to stdin.

	// Some commands will exit when stdin is closed.
	inw.Close()

	// Other commands need a bonk on the head.
	cmd.Process.Signal(os.Interrupt)

	select {
	case <-done:
	case <-time.After(time.Second):
		// A bigger bonk on the head.
		cmd.Process.Signal(os.Kill)
		<-done
	}

	cmd.Process.Wait()
}
```

And, the pump functions:

``` golang
func pumpStdin(ws *websocket.Conn, w io.Writer) {
	defer ws.Close()

	for { // For each incoming message...
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n') // Payload should have a newline at the end.

		// Write headers.
		_, err = fmt.Fprintf(w, "Content-Length: %d\n\n", len(message))
		if err != nil {
			break
		}

		// Write payload.
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	rd := bufio.NewReader(r)
L:
	for { 
		var length int64
		for { // Read headers from stdout until empty line.
			line, err := rd.ReadString('\n')
			if err == io.EOF {
				break L
			}
			if line == "" {
				break
			}
			colon := strings.Index(line, ":")
			if colon < 0 {
				break
			}
			name, value := line[:colon], strings.TrimSpace(line[colon+1:])
			switch name {
			case "Content-Length":
				// Parse Content-Length header value.
				length, _ = strconv.ParseInt(value, 10, 32)
			}
		}

		// Read payload.
		data := make([]byte, length)
		io.ReadFull(rd, data)

		// Write payload to WebSocket.
		if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
			ws.Close()
			break
		}
	}
	close(done)

	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}
```

## And, That's It!


