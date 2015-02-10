// following along http://facebook.github.io/react/docs/tutorial.html

import React from "react";

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
		return (
			<div className="comment-list">
				<Comment author="Alice">Down the... halfpipe!</Comment>
				<Comment author="PP">Alice, what blasphemy!</Comment>
				<Comment author="Maybe Not Lewis">Let her be herself, Peter,
					we let you jump around all the time as well...
				</Comment>
			</div>
		);
	}
}

class CommentBox extends React.Component {
	render() {
		return (
			<div className="comment-box">
				<h1>Comments</h1>
				<CommentList />
				<CommentForm />
			</div>
		);
	}
}

class Comment extends React.Component {
	render() {
		return (
			<div className="comment">
				<h2 className="comment-author">{this.props.author}</h2>
				{this.props.children}
			</div>
		);
	}
}

React.render(<CommentBox />, document.getElementById("content"));
