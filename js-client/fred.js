'use strict';

const request = require("request-promise-native");

module.exports = class Fred {
    constructor(host, port, apiversion) {
        // it would be way cooler to have a template string here, but WebStorm doesn't really let me do that
        // url: `http://{this.host}:{this.port}/{this.baseapi}`
        this.baseurl = 'http://' + this.host + ':' + this.port + '/' + apiversion
    }

    async read(kg, id) {
        let options = {
            url: this.baseurl + '/keygroup/' + kg + '/items/' + id
        };

        return await request(options);
    }

    async put(kg, id, data) {
        const dataString = {
            data: data
        };

        let options = {
            url: this.baseurl + '/keygroup/' + kg + '/items/' + id,
            method: 'PUT',
            body: JSON.stringify(dataString)
        };

        return await request(options);
    }

    async delete(kg, id) {
        let options = {
            url: this.baseurl + '/keygroup/' + kg + '/items/' + id,
            method: 'DELETE'
        };

        return await request(options);
    }

    async createKeygroup(kg) {
        let options = {
            url: this.baseurl + '/keygroup/' + kg,
            method: 'POST'
        };

        return await request(options);
    }

    async deleteKeygroup(kg) {
        let options = {
            url: this.baseurl + '/keygroup/' + kg,
            method: 'DELETE'
        };

        return await request(options);
    }
};