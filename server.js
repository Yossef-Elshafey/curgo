const http = require('http');
const port = 8000;

const resBody = {
  "fname":"Yossef",
  "lname":"Elshafey",
}
http.createServer(async (req, res) => {
  if ( req.url === "/" ) {
    console.log("Got call at", req.url)
    res.writeHead(200, {
      'Content-Type': 'text/plain',
      'cookie': [
        'name=example; Max-Age=9000; HttpOnly',
        'preferences=dark; Expires=Wed, 09 Jun 2021 10:18:14 GMT',
        'sessionToken=abc123; Path=/; Secure; HttpOnly',
        'shoppingCart=12345; Domain=example.com;',
        'logged_in=true; Secure; Path=/; Domain=example.com; HttpOnly'
      ]
    });
    // res.end()
    res.end(JSON.stringify(resBody));
  }

  if (req.url === "/api/test") {
  console.log("Got call at", req.url)
	res.end();
  }
}).listen(port, () => {
  console.log(`App is running on port ${port}`);
  });
