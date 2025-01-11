# School of Computing &mdash; Year 4 Project Proposal Form

## SECTION A

|                     |                              |
|---------------------|------------------------------|
| Project Title:      | GoForMail                    |
| Student 1 Name:     | Andre Rafael Cruz da Fonseca |
| Student 1 ID:       | 21460092                     |
| Student 2 Name:     | Sean Albert Dagohoy          |
| Student 2 ID:       | 21392656                     |
| Project Supervisor: | Stephen Blott                |

## SECTION B

### Introduction

> GoForMail is a mailing list manager which will automatically redirect any received email to the users of a
> configured mailing list. The objective is to empower companies to effortlessly send emails to groups of people
> that may be constantly changing.

### Outline

> GoForMail provides the ability to create, edit, delete and disable mailing lists through a web interface, locked
> behind an admin log in system. The web interface was chosen as it provides a user-friendly solution to managing
> these lists.
>
> The program also provides the ability to 'lock' mailing lists so emails need to be approved by an admin before
> being forwarded. These approvals, as well as other actions are then stored on an audit log describing what actions
> admins have performed.
>
> GoForMail does not provide its own MTA (Mail Transfer Agent) and instead communicates with an already existing one
> using LMTP (Local Mail Transfer Protocol).

### Background

> During a discussion with previous members of the Redbrick committee, GNU Mailman was brought up as a source of
> problems in Redbrick's systems. The issue stems from the program's ambiguous error logging results in difficult
> debugging, wasting sysadmins' time.
>
> To resolve these problems, we decided to rewrite Mailman with more descriptive errors and in a modern language.
> The aim is to provide Redbrick with a better mailing list manager which won't cause them problems in the future.

### Achievements

> GoForMail is aimed at organisations with many members which will require email sending en-masse to certain groups.
> The following features are aimed at ensuring these organisations will have a smoother experience.
>
> #### Mailing list management
> The application will provide users with a way to edit, delete and disable mailing lists that they create. This is
> a quality of life feature which saves time on having to create a new mailing list every time a change needs to be
> made.
>
> #### Admin approval for mail
> In some instances, users may have mailing lists that they do not wish most users to use. For example, if a company
> has a mailing list of all employees for important updates, they might not wish that any employee is able to mail the
> whole company.
>
> To facilitate instances such as this, we provide the ability to 'lock' mailing lists so all emails to the mailing list
> must be approved by an admin before being forwarded to the users on the mailing list.
> 
> #### Audit logging
> The application provides logging for any actions performed by admins, traceable by account username.
> 
> #### User-friendly interface
> The application provides a web panel from which it can be managed. This allows for easier interaction with the app
> from users.

### Justification

> GoForMail enables organisations to send emails en-masse to desired large groups of specific people. This can be
> achieved through the creation of mailing lists which contain a group of users which emails need to be sent to.
> 
> An example of such an organisation is Redbrick, DCU's Computer Networking Society. They send emails to their
> membership weekly, as well as graduates, root holders and committee occasionally.
> 
> This app is useful to them as it allows them to send these emails without having to manually keep track of every
> person who belongs in one of these groups, and instead can just set the list as the receiver.
> 
> In addition, the Go programming language has been rising in popularity in regard to writing microservices. This
> project allows for us to learn this language which will benefit us in the future.

### Programming language(s)

> The backend of the app is written in Go. This includes the forwarding of emails to lists through LMTP.
> 
> The frontend of the app is written in Typescript, using a web development framework (likely Next.js) and compiled
> into static HTML and CSS to be hosted by the Go app.

### Programming tools / Tech stack

> #### GoLand
> GoLand is our main IDE for programming the app. It is used for both the Go and website sections of the app.
> #### Jira
> Jira is how we keep track and divide work throughout development. It aids us in understanding the current progress of
> the application's development.
> #### Gitlab
> Gitlab acts as our software versioning, allowing us to work on different sections of the app at the same time with
> minimal conflicts. It also runs our pipelines to check linting, run our unit tests and check for vulnerabilities
> (TBC).
> #### Snyk
> Snyk is run through our pipeline and checks for vulnerabilities within our imported libraries (TBC).
> #### Postfix
> Our deployment runs Postfix alongside GoForMail. This is due to the app's dependency on a Mail Transfer Agent. Postfix
> was chosen as it is Redbrick's MTA of choice.

### Hardware

> A server is used for deployment and testing of the application.

### Learning Challenges

> List the main new things (technologies, languages, tools, etc) that you will have to learn.
> #### Technologies
> - LMTP protocol
> - Next.js (or similar framework)
> #### Languages
> - Go
> - Typescript
> #### Tools
> - Gitlab CI
> - Snyk
> - Goland
> - Postfix

### Breakdown of work

#### Student 1 - Andre

> Andre will be responsible for setting up the tech stack (CI, Deployment), completing Jira tickets on the Go
> application and reviewing Sean's code.

#### Student 2 - Sean

> Sean will be responsible for creating the web interface, completing Jira tickets on the Go application and reviewing
> Andre's code.

<br>

<p align="center">
  <img src="../res/gfm.png" width="300px" alt="Gopher delivering mail">
</p>

