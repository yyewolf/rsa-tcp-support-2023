# Example transactions

## As a client
```
{"t":0,"d":{"i":1}}
{"t":1,"d":{}}
{"t":2,"d":{"s":true}}
{"t":4,"d":{"c":1}}
{"t":4,"d":{"c":1}}
{"t":4,"d":{"c":1}}
{"t":4,"d":{"c":1}}
{"t":5,"d":{"m":"Bonjour j'ai un soucis avec ma cafetiere"}}
{"t":4,"d":{"c":1}}
{"t":4,"d":{"c":1}}
{"t":4,"d":{"c":1}}
{"t":5,"d":{"m":"Bonjour, quel est le soucis ?","t":"","f":"John"}}
{"t":4,"d":{"c":1}}
{"t":6}
{"t":5,"d":{"m":"You have been elevated to level 1","t":"","f":"SYSTEM"}}
{"t":4,"d":{"c":0}}
{"t":6}
{"t":5,"d":{"m":"You have been elevated to level 2","t":"","f":"SYSTEM"}}
{"t":4,"d":{"c":0}}
{"t":4,"d":{"c":0}}
{"t":4,"d":{"c":0}}
```

## As an agent
```
{"t":0,"d":{"i":2}}
{"t":1,"d":{"a":"password","n":"John"}}
{"t":2,"d":{"s":true}}
{"t":7,"d":{"ids":[]}}
{"t":5,"d":{"m":"Bonjour j'ai un soucis avec ma cafetiere","t":"","f":"1"}}
{"t":5,"d":{"m":"Bonjour, quel est le soucis ?","t":"1"}}  
{"t":8,"d":{"i":1}}
{"t":9,"d":{"m":[{"m":"Bonjour j'ai un soucis avec ma cafetiere","t":"","f":""},{"m":"Bonjour, quel est le soucis ?","t":"","f":"John"}]}}
{"t":6,"d":{"i":1}}
```

# Packets

| Name | Type | Data |
| :-----: | :---: | :---: |
| Hello | 0   | i: ID (int)\*   |
| Identify | 1   | a: Auth (string)\*\*<br>n: Name (string)\*\*|
| Success | 2   | s: Success (bool)   |
| Error | 3   | e: Error (string)   |
| Agent Count\* | 4   | c: Count (int)   |
| Message | 5   | m: Message (string)\*<br>f: From (string)<br>t: To (string)\*\*   |
| Elevate | 6   | i: ID (int)\*\* |
| Client Present | 7   | ids: IDs ([]int)\*\* |
| Client History Request\*\* | 8   | i: ID (int) |
| Client History Response\*\* | 9   | m: Messages ([]Message packets) |


\* : Required/Only for everyone

\*\* : Required/Only for agent

\*\*\* : Required/Only for client