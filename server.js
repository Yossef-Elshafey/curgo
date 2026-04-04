const http = require('http');
const port = process.env.PORT || 8000;

http.createServer((req, res) => {
  if ( req.url === "/" ) {
	console.log(req.headers['user-agent'])
	console.log(req.headers['content-type'])
	console.log(req.method)
	res.writeHead(200, { "Content-Type": "plain/text" });
	res.write("Hello, world!");
	res.end();
  }

  if (req.url === "/api/test") {
	res.write("Done.")
	res.end()
  }
}).listen(port, () => {
	console.log(`App is running on port ${port}`);
  });
