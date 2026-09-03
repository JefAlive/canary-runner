# ⚡ Canary Runner

<p align="center">
  <b>Chega de abrir 5 terminais e quebrar a cabeça com arquivos de configuração só para ligar um servidor de Tibia.</b>
</p>

<p align="center">
  O <b>Canary Runner</b> é um painel de controle visual e moderno feito para você <b>instalar, configurar, ligar, monitorar e jogar</b> o Canary OT Server (com Tibia 15.25) em um único lugar, sem complicações.
</p>

---

## 😫 Você só queria jogar, mas se deparou com isso...

Subir um servidor moderno de Tibia no seu computador (especialmente usando **Windows + WSL**) virou um processo cansativo e confuso:

* Ter que abrir **4 janelas de terminal** diferentes (uma pro Banco MySQL, uma pro Canary, uma pro Login Server e outra pro Patcher).
* **O IP do WSL muda toda vez que o Windows reinicia**, e aí o jogo não conecta mais até você caçar onde editar o IP de novo.
* Ter que abrir arquivos de texto estranhos (`local.toml`, `config.lua`, `config.toml`) e torcer para não errar uma vírgula.
* **Perder progresso e casas (Rollback)** ao fechar a janela do servidor na pressa sem o Canary salvar o mapa.
* Ter que instalar programas pesados de banco de dados só para criar uma conta ou colocar um personagem como **GOD**.

---

## 🪄 Como o Canary Runner resolve tudo isso?

Com um único comando no terminal, você tem uma central completa:

```text
 🚀 Instalação 1-Click       -> Clona os projetos, baixa o Tibia 15.25 e extrai sozinho
 🖥️ Amigo do Windows (WSL)   -> Detecta seu IP sozinho e abre o Tibia.exe na tela do Windows
 🛑 Desligamento Seguro      -> Avisa o Canary para salvar mapa e jogadores antes de fechar
 👥 Gerenciador de Contas    -> Cria contas, personagens GOD e coloca Tibia Coins em segundos
 📜 Logs Coloridos ao Vivo   -> Acompanhe o servidor ligando em tempo real na mesma janela
 🎨 9 Temas Visuais          -> Escolha a paleta de cores que mais combina com seu estilo
```

---

## 🖥️ Espie como é por dentro

### 1. Menu Principal (Dia a dia descomplicado)
Tudo o que você precisa para controlar seu servidor com as setas do teclado:

```text
┌──────────────────────────────────────────────────────────────────────────┐
│ MySQL: ● ON   Login: ● ON   Canary: ● ON                                 │
└──────────────────────────────────────────────────────────────────────────┘

 ▶ 🚀 Canary Server (Servidor do Jogo)                              ● ON
      Liga o mundo do jogo, monstros, mapa, magias e sistemas

    🔑 Login Server (Servidor de Contas)                             ● ON
      Liga o serviço que valida senhas e lista os personagens

    🗄️ MySQL Database (Banco de Dados)                              ● ON
      Liga onde ficam salvas as contas, personagens, itens e casas

    🎮 Abrir Tibia no Windows                                     ⚡ EXECUTAR
      Dispara o jogo no Windows direto pelo terminal sem caçar a pasta

    👥 Gerenciar Jogadores & Contas                               ⚡ SUBMENU
      Painel para criar contas, criar personagens, GOD e Tibia Coins

    ⚙️ Setup & Conexões do Servidor                               ⚡ SUBMENU
      Instalação 1-click, sincronizar IPs, baixar client e pastas

    🎨 Escolher Tema                                               dracula
      Alterne as cores visuais do painel (Dracula, Tokyo, Tron, etc.)

    ❌ Sair

⚡ ATALHOS: [↑/↓] Navegar  •  [ENTER] Selecionar  •  [T] Temas  •  [Q] Sair
```

---

### 2. Acompanhamento de Logs ao Vivo (Com Safe Shutdown)
Ao clicar no Canary ou Login Server, você vê o servidor subindo em tempo real com as cores originais:

```text
🖥️  PAINEL DO SERVIÇO: CANARY OT SERVER   Status: ● ON   (PID: 23590)
╭──────────────────────────────────────────────────────────────────────────╮
│ [INFO] Loading map: canary.otbm... Done (12.4s)                          │
│ [INFO] Spawning monsters and bosses... Done (4.1s)                       │
│ [INFO] Loaded 1420 spells and 320 raids.                                 │
│ [INFO] >> Canary Server is ONLINE! Running on 0.0.0.0:7171 <<            │
╰──────────────────────────────────────────────────────────────────────────╯

⚡ AÇÕES: [S] 🛑 Safe Stop (Salvar) • [K] 💀 Force Kill (-9) • [Esc] Voltar
```

> **Dica:** Apertando **`[S]`**, ele avisa o Canary para salvar todo o mundo e desconectar os jogadores com integridade antes de fechar. Sem sustos de *rollback*!

---

### 3. Painel de Setup & Configurações Automáticas
Seja jogando pelo Windows no WSL, em um servidor na nuvem ou em Linux puro, o menu se adapta ao que você precisa:

```text
┌──────────────────────────────────────────────────────────────────────────┐
│ ⚙️  SETUP & CONEXÕES   [ Modo: 🖥️ Jogando pelo Windows (WSL) ]           │
└──────────────────────────────────────────────────────────────────────────┘

 ▶ ⚡ Instalação Automática Completa (1-Click)
      Baixa todos os arquivos do servidor, baixa o jogo e deixa tudo pronto

   🔄 Sincronizar Conexão do Jogo com o Windows
      Ajusta o servidor e o jogo com o endereço atual para você conseguir entrar

   🎮 Configurar o Tibia (Client.exe)
      Prepara o jogo para abrir e conectar direto no seu servidor

   📥 Apenas Baixar o Tibia 15.25
      Baixa os arquivos do jogo prontos para jogar se você ainda não tiver

   📁 Alterar Pastas dos Arquivos
      Apenas para quem já baixou ou moveu as pastas do servidor manualmente

   🌐 Trocar Modo de Uso (Ativo: 🖥️ Windows/WSL)
      Mude se você estiver rodando em uma máquina na nuvem ou Linux puro

   ❮ Voltar ao Menu Principal
```

---

### 4. Gestão de Contas e GOD (Sem mexer em Banco de Dados)
Chega de instalar softwares pesados para gerenciar seu servidor local:

```text
┌──────────────────────────────────────────────────────────────────────────┐
│ 👥  GERENCIAR JOGADORES & CONTAS                                         │
└──────────────────────────────────────────────────────────────────────────┘

 1. 👤 Criar Nova Conta
    -> Cria uma nova conta de acesso (Login e Senha) direto no banco

 2. 🧙‍♂️ Criar Novo Personagem
    -> Cria um personagem escolhendo Nome, Vocação, Sexo e Nível inicial

 3. ⭐ Promover Personagem para GOD (Administrador)
    -> Dá poderes de administrador a um personagem para usar comandos no jogo

 4. 💰 Adicionar Tibia Coins na Conta
    -> Adiciona moedas na Store do jogo para comprar cosméticos e itens

 5. 📋 Listar Jogadores & Contas
    -> Exibe a lista com todas as contas, personagens criados e cargos

 6. 🗑️ Deletar Personagem ou Conta

 7. ❮ Voltar ao Menu Principal
```

---

## 🎨 Cores e Temas que Não Cansam a Vista

Pressione a tecla **`[T]`** a qualquer momento para alternar entre 9 temas modernos com suporte a **24-bit TrueColor**:

* 🧛 **Dracula** *(Roxo grafite clássico)*
* 🌌 **Tokyo Night** *(Azul escuro noturno)*
* 🌸 **Catppuccin Mocha** *(Tons pastéis confortáveis)*
* 🌊 **Nord Arctic** *(Azul ardósia glacial)*
* 🍂 **Gruvbox Dark** *(Retrô com tons de terra/madeira)*
* 🌊 **Kanagawa Wave** *(Estilo arte tradicional japonesa)*
* 💠 **Tron: Legacy** *(Azul néon futurista)*
* 🔴 **Tron: Ares** *(Vermelho obsidiana)*
* 📟 **Matrix** *(Verde terminal cyberpunk)*

*(O tema escolhido fica salvo para as próximas vezes que você abrir o painel!)*

---

## 🚀 Como Usar na sua Máquina

### 1. Requisitos Básicos
No seu terminal (Ubuntu no WSL ou Linux nativo), instale os pacotes básicos:

```bash
sudo apt update
sudo apt install -y golang-go git curl unzip mysql-server
```

### 2. Baixar e Rodar
```bash
# Clone este repositório
git clone https://github.com/seu-usuario/canary-runner.git
cd canary-runner

# Inicie o painel!
go run ./src
```

### 3. (Opcional) Abrir de qualquer lugar do terminal
Se quiser digitar apenas `canary-runner` em qualquer pasta:

```bash
go build -o canary-runner ./src
sudo mv canary-runner /usr/local/bin/

# Pronto! Agora basta digitar:
canary-runner
```

---

## ⌨️ Comandos do Teclado

* **`↑ / ↓`** ou **`k / j`**: Navega entre as opções
* **`ENTER`**: Entra no menu, liga o servidor ou executa a ação
* **`S`**: *Safe Stop* (Salva o mapa com segurança e desliga o servidor)
* **`K`**: *Force Kill* (Encerramento forçado caso o servidor trave)
* **`T`**: Abre o seletor de temas visuais
* **`R`**: Atualiza a tela de logs
* **`ESC`** ou **`B`**: Volta para o menu anterior
* **`Q`**: Fecha o painel *(o servidor continua rodando em segundo plano)*

---

<p align="center">
  Feito para simplificar a vida da comunidade <b>OpenTibia & Canary</b> ⚡
</p>