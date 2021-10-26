// send the provided request to api server
export const Request =
    (
        uri,
        method,
        jwt = false,
        jsonBody = false,
        params = false,
        signal = null,
        isMultiPartFormData = false,
        multiPartContentTypeString = "",
        isPaginate = true,
    ) => {
        // prepare request headers
        let reqHeaders = {

        }
        if (isMultiPartFormData) {
            reqHeaders = {
                // Accept: 'application/json',
                // 'Content-Type': 'multipart/form-data; boundary=----WebKitFormBoundaryIn312MOjBWdkffIM',
            };
        } else {
            reqHeaders = {
                Accept: 'application/json',
                'Content-Type': 'application/json',
            };
        }

        // check if has jwt and if yes, include it
        if (jwt) {
            reqHeaders['Authorization'] = 'Bearer ' + jwt;
        }


        // preparing request header
        let reqInit = {
            // Request method
            method: method,

            // Request headers
            headers: { ...reqHeaders, "FromPage": window.location },
        };

        var _uri = uri;

        // console.log(params, typeof params);

        // add params to uri if needed
        if (params !== false && typeof params === 'object') {
            const keys = Object.keys(params);
            keys.map((k, i) => {
                if (i === 0) {
                    _uri += '?';
                } else {
                    _uri += '&';
                }
                _uri += k + '=' + params[k];
                return _uri
            });
        }

        if (jsonBody !== false && typeof jsonBody === 'object') {
            if (isMultiPartFormData) {
                // form data
                const formData = new FormData();
                for (const name in jsonBody) {
                    formData.append(name, jsonBody[name]);
                }
                reqInit['body'] = formData;
            } else {
                // json
                reqInit['body'] = JSON.stringify(jsonBody);
            }
        }


        // console.log(reqInit);

        return (fetch(_uri, reqInit)
            // check response status
            .then((response) => {
                // http response status
                const stat = response.status;

                let jsonResp = {
                    ok: response.ok,
                    status: stat,
                    statusText: response.statusText,
                    i18nMessage: 'http.error.' + stat + '.msg',
                    response: response.json(),
                    error: false,
                };



                // Informational responses
                if (stat >= 100 && stat < 200) {
                    jsonResp.error = false;
                }
                // Successful responses
                else if (stat >= 200 && stat < 300) {
                    jsonResp.error = false;
                }
                // Redirect
                else if (stat >= 300 && stat < 400) {
                    jsonResp.error = false;
                }
                // Client error
                else if (stat >= 400 && stat < 500) {
                    jsonResp.error = true;
                    if (stat === 400) { // input error
                        jsonResp.i18nMessage = 'http.error.' + stat + '.msg'
                    }
                }
                // Server error
                else if (stat >= 500 && stat < 600) {
                    jsonResp.error = true;
                } else {
                    // Nothing here :))
                }

                // return jsonResponse that we have made
                return jsonResp;
            })
            .then((json) => {
                return {
                    ...json,
                };
            })
            .catch((e) => {
                return {
                    ok: false,
                    status: -1,
                    statusText: e,
                    i18nMessage: 'http.resp.unknown',
                    response: e,
                    error: true,
                };
            }));
    };