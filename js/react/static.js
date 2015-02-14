import fs from "fs";
import React from "react";

import {CommentBox} from "./comments.js";

let commentsHtml = React.renderToString(<CommentBox url="comments.json" />);

let html = `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Hi, React! (static...)</title>
  </head>

  <body>
    <div id="content">
    	${commentsHtml}
    </div>

    <script src="bundle.js"></script>
  </body>
</html>
`

console.log(html);
