FROM centurylink/ca-certs
COPY ./dist/linux/dsfs /dsfs
EXPOSE 8040
ENTRYPOINT ["/dsfs"]
