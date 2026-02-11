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

### How to verify you won the challenge

1. Run your server.
2. Open Chrome DevTools -> **Network** tab.
3. Refresh the page.
4. Look at the `localhost` (or `index`) request.
5. Check the **"Transfer Size"** (this is the compressed size sent over the wire).
6. If it is **< 14.6 kB**, you win!

Would you like to start building this? I can guide you through setting up the `go.mod` and running the generation steps!
