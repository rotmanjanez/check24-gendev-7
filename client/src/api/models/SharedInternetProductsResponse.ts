/**
 * CHECK24 GenDev 7 API
 * API for the 7th CHECK24 GenDev challenge providing product offerings from five different internet providers
 *
 * OpenAPI spec version: dev
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { Address } from '../models/Address';
import { InternetProduct } from '../models/InternetProduct';

/**
* Response containing a list of shared internet products
*/
export class SharedInternetProductsResponse {
    'products'?: Array<InternetProduct>;
    'version'?: string;
    'address'?: Address;

    static readonly discriminator: string | undefined = undefined;

    static readonly mapping: { [index: string]: string } | undefined = undefined;

    static readonly attributeTypeMap: Array<{ name: string, baseName: string, type: string, format: string }> = [
        {
            "name": "products",
            "baseName": "products",
            "type": "Array<InternetProduct>",
            "format": ""
        },
        {
            "name": "version",
            "baseName": "version",
            "type": "string",
            "format": ""
        },
        {
            "name": "address",
            "baseName": "Address",
            "type": "Address",
            "format": ""
        }];

    static getAttributeTypeMap() {
        return SharedInternetProductsResponse.attributeTypeMap;
    }

    public constructor() {
    }
}
