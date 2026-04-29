# Installation

# DEV
- [x] Tokens position implemention for better error message
- [x] Eval FetchStatement function needs to be implemented
- [x] There has to be a seperation between expectPeekToken function and the error
- [x] (search(noPrefixFoundErr)) function has to be implemented
- [x] Eval if else
- [ ] create an api between object.Response and curgo, after that
      response should not abort the evaluator
      added:
      1- the api should have a return value for example if you want to store a value of a header
      how would youa access it ? it should be a hash table where you could access values like x[res]
      2- object.Response has to implemented in a such a way that accessing its
      props and method make sense in terms of programming behaviors

- [ ] implement hash maps to be the interface of map in host, to be used in response accessing
- [ ] there has to be a concrete way to handle errors in evaluator since ignoring the check of
      an object.Error(isError) make it keep flying while tree walking
