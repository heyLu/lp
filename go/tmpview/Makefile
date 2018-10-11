NAME = tmpview

$(NAME): $(NAME).go
	go build .

install: $(NAME)
	install -m0755 $(NAME) $(PREFIX)/bin
