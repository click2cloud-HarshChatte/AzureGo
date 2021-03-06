# Create a Base Image     
FROM golang:latest 
# Setting up the Working Directory     
WORKDIR /app       
# Git Installer     
RUN apt-get update && apt-get install git       
# Setting up Run Files             
RUN go get github.com/gorilla/mux  
RUN go get github.com/Azure/go-autorest/autorest/to      
RUN go get github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-11-01/network  
RUN go get github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources       
RUN go get github.com/Azure/go-autorest/autorest/azure/auth        
#copying the files into the working directory                    
COPY /AzureGo /app                    
# Exposing or defining port to access to outside world          
EXPOSE 8080                                                    
# Building  the main file of go                             
RUN go build -o main .                              
# Running the main file           
CMD [ "/app/main" ]
