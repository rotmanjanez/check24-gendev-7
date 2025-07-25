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

export class Discount {
    /**
    * Determines which subtype to use
    */
    'type': string;

    static readonly discriminator: string | undefined = "type";

    static readonly mapping: { [index: string]: string } | undefined = {
        "absolute": "AbsoluteDiscount",
        "percentage": "PercentageDiscount",
    };

    static readonly attributeTypeMap: Array<{ name: string, baseName: string, type: string, format: string }> = [
        {
            "name": "type",
            "baseName": "type",
            "type": "string",
            "format": ""
        }];

    static getAttributeTypeMap() {
        return Discount.attributeTypeMap;
    }

    public constructor() {
        this.type = "Discount";
    }
}
