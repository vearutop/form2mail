FROM alpine

COPY ./bin/form2mail /bin/form2mail

EXPOSE 80

ENTRYPOINT [ "/bin/form2mail" ]
