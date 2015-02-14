// following along http://facebook.github.io/react/docs/tutorial.html

import React from "react";
import {CommentBox} from "./comments.js";

require('./style.css');

React.render(<CommentBox url="comments.json" />, document.getElementById("content"));
