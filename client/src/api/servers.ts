import { RequestContext, HttpMethod } from "./http/http";

export interface BaseServerConfiguration {
    makeRequestContext(endpoint: string, httpMethod: HttpMethod): RequestContext;
}

/**
 *
 * Represents the configuration of a server including its
 * url template and variable configuration based on the url.
 *
 */
export class ServerConfiguration<T extends { [key: string]: string }> implements BaseServerConfiguration {
    public constructor(private url: string, private variableConfiguration: T) { }

    /**
     * Sets the value of the variables of this server. Variables are included in 
     * the `url` of this ServerConfiguration in the form `{variableName}`
     *
     * @param variableConfiguration a partial variable configuration for the 
     * variables contained in the url
     */
    public setVariables(variableConfiguration: Partial<T>) {
        Object.assign(this.variableConfiguration, variableConfiguration);
    }

    public getConfiguration(): T {
        return this.variableConfiguration
    }

    private getUrl() {
        let replacedUrl = this.url;
        for (const [key, value] of Object.entries(this.variableConfiguration)) {
            replacedUrl = replacedUrl.replace(new RegExp(`\\{${key}\\}`, 'g'), value);
        }
        return replacedUrl
    }

    /**
     * Creates a new request context for this server using the url with variables
     * replaced with their respective values and the endpoint of the request appended.
     *
     * @param endpoint the endpoint to be queried on the server
     * @param httpMethod httpMethod to be used
     *
     */
    public makeRequestContext(endpoint: string, httpMethod: HttpMethod): RequestContext {
        return new RequestContext(this.getUrl() + endpoint, httpMethod);
    }
}

// no server specified so that vite can proxy requests to the backend
export const dev = new ServerConfiguration<{}>("/api", {})
export const prod = new ServerConfiguration<{}>("https://check24-gendev-7-api.janez.at/api", {})

export const servers = [prod, dev];
