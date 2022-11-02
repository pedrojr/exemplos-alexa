
Criar nova skill
* Clicar em "Create Skill"
* Adicionar um nome para "ImplicitGrantLogin"
* Selecionar o modelo "Custom"
* Selecionar a opção "Alexa-hosted (Python)"
* Clicar em "Create Skill"
* Aguardar...
* Escolher a opção "Start from Scratch"
* Clicar em "Continue with template"
* Aguardar...
* Na aba "Build"
* Na opção "Invocations" > "Skill Invocation Name" altere o nome da skill para "listar itens"
* Clicar em Salvar
* Na opção "Interaction Model" > "Intents" editar a intent "HelloWorldIntent"
	* Alterar nome para ListarItensIntent
	* Remover enunciados
	* Adicionar o enunciado "listar"
* Clicar em "Build Model"

---

Configurar o "Account Linking" com "Implicit Grant"
* [Configurar implicit grant](../02-Amazon-Implicit-Grant.md)

---

Atualizar código da skill
* Na aba "Code", abrir o arquivo "lambda_function.py" e substituir pelo código do repositório "/skill/lambda_function.py"
* Ajustar a url no código da skill
