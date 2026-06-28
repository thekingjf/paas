# Design Decisions

## How to use this file
Log an entry only if it clears the bar:
- DECISION: a competent engineer could have reasonably chosen otherwise.
Keep entries 2–4 lines. Capture the WHY / the lesson, not the code.
When unsure, log it — under-logging (forgetting a good story) is the bigger risk.

---
## Template
### <short title of the decision>
- **Fork:** <the choice you faced, one line>
- **Chose:** <what you picked / what you rejected>
- **Why:** <the reasoning — the defensible part>

### Driver Choice
- **Fork:** Had to decide between modernc vs mattn for my SQLite drivers for GO
- **Chose:** I chose modernc
- **Why:** Modernc simplifies deployment and allows future me to have an easier time
           when it comes up and as a personal project I don't need the best I need
           something that works consistently and easily without becoming it's own
           issue
    
### Server Struct Design
- **Fork:** Handlers needed access to the DB and 3 valid options:
            - Package level global DB var
            - pass db into every handler
            - Hang it onto a struct and with handler methods
- **Chose:** I chose to hang it on a struct with handler methods
- **Why:** A struct scales cleanly as dependcies increase and make handlers testable

### Parameterized Queries
- **Fork:** Fmt.sprintf vs ? placeholders for the sql string
- **Chose:** I chose to move forward with ? placeholders
- **Why:** I chose to do the ? placeholders for security because it passes the
           value over to the driver separately meaning it's never processed as SQL,
           preventing SQL injection that sprintf would have

### No Nulls Default Fill
- **Fork:** Inserting basic placeholder values vs leaving them NULL
- **Chose:** I chose to insert basic placeholder values("", 0, "created")
- **Why:** I chose to insert placeholders because a created row in this databse 
           implies existence so the placeholder values are geninune and not covering  
           for missing data.

### Errors.Is vs Errors.As
- **Fork:** Recognizing which to use
- **Chose:** I chose to use Errors is
- **Why:** When I need to check an error I know it's between if I need to match
           it(Is) or if I need to extract it(As)

### Docker Migration
- **Fork:** stay on docker/docker@v28(old) or migrate moby/moby/client
- **Chose:** I chose to migrate to the newer docker version instead of staying
             on the old frozen version
- **Why:** I chose to migrate to the newer version because it's reflective of the
           current state of docker which is most likely actually used by the 
           majority of people and since it's early on I won't need to make deep 
           changes and can avoid potential bugs that the old version help that were
           fixed in the newer version before i ever run into them