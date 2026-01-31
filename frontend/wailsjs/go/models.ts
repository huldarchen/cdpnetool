export namespace api {
	
	export class Response_cdpnetool_internal_gui_BrowserData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.BrowserData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_BrowserData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.BrowserData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_ConfigData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.ConfigData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_ConfigData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.ConfigData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_ConfigListData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.ConfigListData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_ConfigListData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.ConfigListData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_EventHistoryData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.EventHistoryData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_EventHistoryData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.EventHistoryData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_NewConfigData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.NewConfigData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_NewConfigData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.NewConfigData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_NewRuleData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.NewRuleData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_NewRuleData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.NewRuleData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_SessionData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.SessionData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_SessionData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.SessionData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_SettingData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.SettingData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_SettingData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.SettingData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_SettingsData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.SettingsData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_SettingsData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.SettingsData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_StatsData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.StatsData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_StatsData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.StatsData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_TargetListData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.TargetListData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_TargetListData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.TargetListData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Response_cdpnetool_internal_gui_VersionData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    data?: gui.VersionData;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_internal_gui_VersionData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], gui.VersionData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EmptyData {
	
	
	    static createFrom(source: any = {}) {
	        return new EmptyData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Response_cdpnetool_pkg_api_EmptyData_ {
	    success: boolean;
	    code?: string;
	    message?: string;
	    // Go type: EmptyData
	    data?: any;
	
	    static createFrom(source: any = {}) {
	        return new Response_cdpnetool_pkg_api_EmptyData_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.code = source["code"];
	        this.message = source["message"];
	        this.data = this.convertValues(source["data"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace domain {
	
	export class EngineStats {
	    total: number;
	    matched: number;
	    byRule: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new EngineStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.matched = source["matched"];
	        this.byRule = source["byRule"];
	    }
	}
	export class TargetInfo {
	    id: string;
	    type: string;
	    url: string;
	    title: string;
	    isCurrent: boolean;
	
	    static createFrom(source: any = {}) {
	        return new TargetInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.type = source["type"];
	        this.url = source["url"];
	        this.title = source["title"];
	        this.isCurrent = source["isCurrent"];
	    }
	}

}

export namespace gui {
	
	export class BrowserData {
	    devToolsUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new BrowserData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.devToolsUrl = source["devToolsUrl"];
	    }
	}
	export class ConfigData {
	    config?: model.ConfigRecord;
	
	    static createFrom(source: any = {}) {
	        return new ConfigData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config = this.convertValues(source["config"], model.ConfigRecord);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ConfigListData {
	    configs: model.ConfigRecord[];
	
	    static createFrom(source: any = {}) {
	        return new ConfigListData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.configs = this.convertValues(source["configs"], model.ConfigRecord);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EventHistoryData {
	    events: model.NetworkEventRecord[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new EventHistoryData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.events = this.convertValues(source["events"], model.NetworkEventRecord);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NewConfigData {
	    config?: model.ConfigRecord;
	    configJson: string;
	
	    static createFrom(source: any = {}) {
	        return new NewConfigData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config = this.convertValues(source["config"], model.ConfigRecord);
	        this.configJson = source["configJson"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NewRuleData {
	    ruleJson: string;
	
	    static createFrom(source: any = {}) {
	        return new NewRuleData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ruleJson = source["ruleJson"];
	    }
	}
	export class SessionData {
	    sessionId: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sessionId = source["sessionId"];
	    }
	}
	export class SettingData {
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new SettingData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.value = source["value"];
	    }
	}
	export class SettingsData {
	    settings: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new SettingsData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.settings = source["settings"];
	    }
	}
	export class StatsData {
	    stats: domain.EngineStats;
	
	    static createFrom(source: any = {}) {
	        return new StatsData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stats = this.convertValues(source["stats"], domain.EngineStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class TargetListData {
	    targets: domain.TargetInfo[];
	
	    static createFrom(source: any = {}) {
	        return new TargetListData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.targets = this.convertValues(source["targets"], domain.TargetInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class VersionData {
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	    }
	}

}

export namespace model {
	
	export class ConfigRecord {
	    id: number;
	    configId: string;
	    name: string;
	    version: string;
	    configJson: string;
	    isActive: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new ConfigRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.configId = source["configId"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.configJson = source["configJson"];
	        this.isActive = source["isActive"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NetworkEventRecord {
	    id: number;
	    sessionId: string;
	    targetId: string;
	    url: string;
	    method: string;
	    statusCode: number;
	    finalResult: string;
	    matchedRulesJson: string;
	    requestJson: string;
	    responseJson: string;
	    timestamp: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new NetworkEventRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sessionId = source["sessionId"];
	        this.targetId = source["targetId"];
	        this.url = source["url"];
	        this.method = source["method"];
	        this.statusCode = source["statusCode"];
	        this.finalResult = source["finalResult"];
	        this.matchedRulesJson = source["matchedRulesJson"];
	        this.requestJson = source["requestJson"];
	        this.responseJson = source["responseJson"];
	        this.timestamp = source["timestamp"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

