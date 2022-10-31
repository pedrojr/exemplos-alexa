
Conta amazon
* Link: https://developer.amazon.com/alexa/console/ask
* Criar uma nova skill usando python
* Na opção "TOOLS" clique em "Account Linking"
* Ativar o account linking
    * Do you allow users to create an account or link to an existing account with you? [ON]
    * Settings > Allow users to enable skill without account linking (Recommended). [ON]
* Na opção "Implicit Grant" configure os campos
    * Em "Your Web Authorization URI" para https://SEU_HOSTNAME.ddns.net:8090/
    * Em "Your Client ID" adicionar o ClientID do seu arquivo .env
    * Em "Scope" adicionar all
    * Em "Domain List" adicionar SEU_HOSTNAME.ddns.net
