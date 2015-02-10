// following along http://facebook.github.io/react/docs/tutorial.html

import React from "react";
import marked from "marked";

require('./style.css');

class CommentForm extends React.Component {
	constructor() {
		this.handleSubmit = this.handleSubmit.bind(this);
	}

	handleSubmit(ev) {
		ev.preventDefault();
		let author = this.refs.author.getDOMNode().value.trim();
		let text = this.refs.text.getDOMNode().value.trim();
		if (!text || !author) {
			return;
		}

		this.props.onCommentSubmit({author, text});
		this.refs.author.getDOMNode().value = '';
		this.refs.text.getDOMNode().value = '';
	}

	render() {
		return (
			<form className="comment-form" onSubmit={this.handleSubmit}>
				<input type="text" placeholder="Your name" ref="author" />
				<input type="text" placeholder="Say something..." ref="text" />
				<input type="submit" value="Say it!" />
			</form>
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
	constructor(props) {
		super(props);
		this.state = {data: []}

		this.handleCommentSubmit = this.handleCommentSubmit.bind(this);
	}

	handleCommentSubmit(comment) {
		let comments = this.state.data;
		let newComments = comments.concat([comment]);
		this.setState({data: newComments});

		// TODO: Send to server
	}

	componentDidMount() {
		let xhr = new XMLHttpRequest();
		xhr.open('GET', this.props.url);
		xhr.onreadystatechange = function() {
			if (xhr.readyState == XMLHttpRequest.DONE) {
				this.setState({data: JSON.parse(xhr.responseText)});
			}
		}.bind(this);
		xhr.onerror = function(err) {
			console.error(this.props.url, xhr.status, err);
		}.bind(this);
		xhr.send();
	}

	render() {
		return (
			<div className="comment-box">
				<h1>Comments</h1>
				<CommentList data={this.state.data} />
				<CommentForm onCommentSubmit={this.handleCommentSubmit} />
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

React.render(<CommentBox url="comments.json" />, document.getElementById("content"));
