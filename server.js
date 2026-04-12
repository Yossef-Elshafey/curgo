const http = require('http');
const port = 8000;

http.createServer(async (req, res) => {
  if ( req.url === "/" ) {
  console.log("Got call at", req.url)
	res.end();
  }

  if (req.url === "/api/test") {
  console.log("Got call at", req.url)
	res.end();
  }
}).listen(port, () => {
  console.log(`App is running on port ${port}`);
  });

