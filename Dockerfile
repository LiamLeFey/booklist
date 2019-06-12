
# eventually, I'd prefer to use a multi-stage process to compile in a golang
# container, and then run in a scratch container.
FROM golang

# Set the working directory to /app
WORKDIR /app

COPY . /app

# Make port 8080 available to the world outside this container
EXPOSE 8080

# Define environment variable
ENV NAME booklistContainer

#Next stage. test building an app
CMD ["./booklist"]

