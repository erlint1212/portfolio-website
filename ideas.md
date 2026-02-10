Make website under 14kb, [primes video](https://www.youtube.com/watch?v=ciNXbR5wvhU)

---

Hello! As your Coding Partner, I'm excited to help you tackle this challenge. The video you shared discusses a fascinating performance optimization principle rooted in how the internet's transmission protocols work.

Here is the breakdown of the **"14kB Challenge"** rules based on the video's principles (TCP Slow Start), tailored for your Go + HTMX + Templ stack.

### The Core Principle: TCP Slow Start

When a connection is established, the TCP protocol doesn't know the network's capacity. It starts by sending a small amount of data (usually **10 packets**, roughly **14kB**) and waits for an acknowledgement (ACK). If your website fits in this initial burst, it loads instantly, saving an entire round-trip of latency (which can be 600ms+ on bad networks).

---

### 🚀 The Challenge Rules

1. **The 14kB Hard Limit (Compressed)**
* **Rule:** Your initial HTTP response (the HTML for `index.html`) must be under **14kB** after Gzip/Brotli compression.
* **Why:** This ensures the browser gets the entire page in the first TCP round trip.


2. **Server-Side Rendering (SSR) is Mandatory**
* **Rule:** The initial page must be fully rendered HTML sent from Go. No "Loading..." spinners that wait for a second JavaScript fetch.
* **Implementation:** Use **Templ** to generate the HTML string on the server.


3. **Zero-Bundle / Minimal Scripting**
* **Rule:** Do not serve a massive `bundle.js`.
* **Implementation:** Use **HTMX** (which is ~14kB gzipped itself, fitting the budget if inline, or cached separately) for interactions.


4. **Inline "Above the Fold" CSS**
* **Rule:** To avoid an extra network request for a `.css` file that blocks rendering, inline the critical CSS directly into the `<head>` of your HTML.
* **Tip:** If you use Tailwind, extract only the classes used on the landing page and put them in a `<style>` tag.


5. **Lazy Load Heavy Assets**
* **Rule:** Images and secondary content (below the fold) must not delay the initial render.
* **Implementation:** Use `loading="lazy"` for images and HTMX's `hx-trigger="reveal"` to load lower sections of the page only when the user scrolls to them.



---

### 🛠️ The Implementation Plan (Go + HTMX + Templ)

Here is a starter scaffolding designed to strictly follow these rules.

#### 1. Project Structure

```text
/cmd/server/main.go   # Go server with Gzip middleware
/views/index.templ    # Templ components
/public/              # Static files (images, full css)
go.mod

```

#### 2. The Code

**Step 1: Install Templ and HTMX**
First, ensure you have the `templ` CLI installed (`go install github.com/a-h/templ/cmd/templ@latest`).

**Step 2: Create the Templ Component (`views/index.templ`)**
We will keep this lightweight. Notice we are using a CDN for HTMX, but ideally, you would host a minified version or inline it if you are strict about the "single packet" concept.

```go
package views

// Standard HTML5 boilerplate
templ Page(title string) {
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8"/>
		<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
		<title>{ title }</title>
		
		<style>
			body { font-family: system-ui, sans-serif; margin: 0; padding: 1rem; line-height: 1.5; }
			.container { max-width: 600px; margin: 0 auto; }
			.btn { background: #007bff; color: white; padding: 0.5rem 1rem; border: none; cursor: pointer; }
		</style>

		<script src="https://unpkg.com/htmx.org@1.9.10" integrity="sha384-D1Kt99CQMDuVetoL1lrYwg5t+9QdHe7NLX/SoJYkXDFfX37iInKRy5xLSi8nO7UC" crossorigin="anonymous"></script>
	</head>
	<body>
		<div class="container">
			<h1>The 14kB Challenge</h1>
			<p>
				This content was rendered on the server. It arrived in the first TCP window.
			</p>
			
			<div hx-get="/heavy-content" hx-trigger="reveal">
				(Scroll down to load more content...)
			</div>
		</div>
	</body>
	</html>
}

templ HeavyContent() {
    <div style="margin-top: 50vh; padding: 2rem; background: #f0f0f0;">
        <h2>Loaded Lazily!</h2>
        <p>This content wasn't in the initial 14kB payload.</p>
    </div>
}

```

**Step 3: Create the Go Server with Compression (`cmd/server/main.go`)**
Compression is *vital* to hit the 14kB target.

```go
package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"your-module-name/views" // Import your generated views
)

// Gzip Middleware to ensure we fit in the packet
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		next.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}

// Wrapper for gzip response
type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func main() {
	// Handlers
	http.Handle("/", gzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Render the main page
		views.Page("14kB Challenge").Render(r.Context(), w)
	})))

	http.Handle("/heavy-content", gzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Render the lazy-loaded content
		views.HeavyContent().Render(r.Context(), w)
	})))

	println("Server starting on :8080 (Check payload size in DevTools!)")
	http.ListenAndServe(":8080", nil)
}

```

### How to verify you won the challenge

1. Run your server.
2. Open Chrome DevTools -> **Network** tab.
3. Refresh the page.
4. Look at the `localhost` (or `index`) request.
5. Check the **"Transfer Size"** (this is the compressed size sent over the wire).
6. If it is **< 14.6 kB**, you win!

Would you like to start building this? I can guide you through setting up the `go.mod` and running the generation steps!
