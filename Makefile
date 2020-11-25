PARSER	=	parser
FILES	=	handlers.go		main.go		model.go	parsing.go	config.go	database.go

.PHONY: all clean fclean re

all:		$(PARSER)

parser:	$(FILES)
			go build -o $(PARSER) $(FILES) 

clean:
			@echo "clean"

fclean:	clean
			@rm -rf $(PARSER)

re:		fclean all