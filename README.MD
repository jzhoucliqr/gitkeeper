run keeper:  

./gitkeeper -t "github-user-token-for this gitkeeper user"


run webhook: 

./webhook -p 8888 -s "SECRETFORHOOK"

example json of repoconfig to send to api to add to gitkeeper:

repo.json:  {"Repo":{"Name":"CloudCenter","FullName":"jzhoucliqr/CloudCenter","WebhookSecret":"SECRETFORHOOK"},"User":{"Name":"jzhoucliqr","Token":"4ac35241943b4931cf6eaee2553"}}

curl -v -X POST -d @repo.json http://127.0.0.1:9999/api 

Need etcd as backend
