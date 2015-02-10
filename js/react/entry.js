// following along http://facebook.github.io/react/docs/tutorial.html

import React from "react";
import marked from "marked";

require('./style.css');

class CommentForm extends React.Component {
	render() {
		return (
			<div className="comment-form">
				I might be an actual form someday...
			</div>
		);
	}
}

class CommentList extends React.Component {
	render() {
		let commentNodes = this.props.data.map((comment) => {
			return (
				<Comment author={comment.author}>
					{comment.text}
				</Comment>
			);
		});

		return (
			<div className="comment-list">
				{commentNodes}
			</div>
		);
	}
}

class CommentBox extends React.Component {
	render() {
		return (
			<div className="comment-box">
				<h1>Comments</h1>
				<CommentList data={this.props.data} />
				<CommentForm />
			</div>
		);
	}
}

class Comment extends React.Component {
	render() {
		let rawMarkup = marked(this.props.children.toString());
		return (
			<div className="comment">
				<h2 className="comment-author">{this.props.author}</h2>
				<span dangerouslySetInnerHTML={{__html: rawMarkup}} />
			</div>
		);
	}
}

var data = [
	{author: "Alice", text: "Down the... halfpipe!"},
	{author: "PP", text: "Alice, what **blasphemy**!"},
	{author: "Maybe Not Lewis", text: "Let her be herself, Peter, we let *you* jump around all the time as well..."}
];

React.render(<CommentBox data={data} />, document.getElementById("content"));
