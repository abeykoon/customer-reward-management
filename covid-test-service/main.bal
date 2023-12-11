import ballerina/http;
import ballerina/url;
import ballerina/os;
import ballerina/log;

service /info on new http:Listener(9090) {
    resource function get stats(string country) returns json|error {

        string clientId = os:getEnv("CLIENT_ID");
        string clientSecret = os:getEnv("CLIENT_SECRET");
        string tokenUrl = os:getEnv("TOKEN_URL");
        string serviceUrl = os:getEnv("SERVICE_URL");

        log:printInfo("ClientId = " + clientId);
        log:printInfo("Client Seret = " + clientSecret);
        log:printInfo("TokenUrl = " + tokenUrl);
        log:printInfo("ServiceUrl = " + serviceUrl);

        http:Client clientEp = check new (
            url = serviceUrl,
            config = {
                auth: {
                    clientId: clientId,
                    clientSecret: clientSecret,
                    clientConfig: {},
                    tokenUrl: tokenUrl
                }
            }

        );
        string encodedCountry = check url:encode(country, "UTF-8");
        string resourcePAth = "/v3/covid-19/countries/" + encodedCountry + "?strict=true";
        http:Response res = check clientEp->get(resourcePAth);
        json response = check res.getJsonPayload();
        log:printInfo(response.toString());
        return response;
    }
}
