'use strict';

const Fred = require('./fred.js');

let fred = new Fred('localhost', '9001');

(async () => {
    console.log(await fred.createKeygroup('kg'));
    console.log(await fred.put('kg', '1', 'hi!'));
    console.log(await fred.read('kg', '1'));
    console.log(await fred.delete('kg', '1'));
    console.log(await fred.deleteKeygroup('kg'));
})();
