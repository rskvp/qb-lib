/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./src/launcher.ts");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./src/collections/Dictionary.ts":
/*!***************************************!*\
  !*** ./src/collections/Dictionary.ts ***!
  \***************************************/
/*! exports provided: Dictionary */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "Dictionary", function() { return Dictionary; });
var Dictionary = /** @class */ (function () {
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function Dictionary(o) {
        // ------------------------------------------------------------------------
        //                      f i e l d s
        // ------------------------------------------------------------------------
        this._items = {};
        this._count = 0;
        if (!!o) {
            for (var key in o) {
                if (o.hasOwnProperty(key)) {
                    this.put(key, o[key]);
                }
            }
        }
    }
    // ------------------------------------------------------------------------
    //                      p u b l i c
    // ------------------------------------------------------------------------
    Dictionary.prototype.putAll = function (data) {
        var items;
        if (data instanceof Dictionary) {
            items = data._items;
        }
        else {
            items = data;
        }
        for (var key in items) {
            if (items.hasOwnProperty(key)) {
                this.put(key, items[key]);
            }
        }
    };
    Dictionary.prototype.put = function (key, value) {
        this._items[key] = value;
        this._count++;
    };
    Dictionary.prototype.get = function (key) {
        return this._items[key];
    };
    Dictionary.prototype.containsKey = function (key) {
        return this._items.hasOwnProperty(key);
    };
    Dictionary.prototype.count = function () {
        return this._count;
    };
    Dictionary.prototype.isEmpty = function () {
        return this._count === 0;
    };
    Dictionary.prototype.keys = function () {
        var Keys = [];
        // tslint:disable-next-line:forin
        for (var key in this._items) {
            Keys.push(key);
        }
        return Keys;
    };
    Dictionary.prototype.remove = function (key) {
        var val = this._items[key];
        delete this._items[key];
        this._count--;
        return val;
    };
    Dictionary.prototype.values = function () {
        var values = [];
        // tslint:disable-next-line:forin
        for (var key in this._items) {
            values.push(this._items[key]);
        }
        return values;
    };
    Dictionary.prototype.clear = function () {
        if (!this.isEmpty()) {
            this._items = {};
        }
    };
    return Dictionary;
}());



/***/ }),

/***/ "./src/commons/BaseObject.ts":
/*!***********************************!*\
  !*** ./src/commons/BaseObject.ts ***!
  \***********************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _random__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./random */ "./src/commons/random.ts");

var BaseObject = /** @class */ (function () {
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function BaseObject() {
        this._uid = _random__WEBPACK_IMPORTED_MODULE_0__["default"].uniqueId(BaseObject.PREFIX);
    }
    Object.defineProperty(BaseObject.prototype, "uid", {
        // ------------------------------------------------------------------------
        //                      p u b l i c
        // ------------------------------------------------------------------------
        get: function () {
            return this._uid;
        },
        enumerable: true,
        configurable: true
    });
    // ------------------------------------------------------------------------
    //                      c o n s t
    // ------------------------------------------------------------------------
    BaseObject.PREFIX = "lyts_object_";
    return BaseObject;
}());
// ------------------------------------------------------------------------
//                      e x p o r t s
// ------------------------------------------------------------------------
/* harmony default export */ __webpack_exports__["default"] = (BaseObject);


/***/ }),

/***/ "./src/commons/console.ts":
/*!********************************!*\
  !*** ./src/commons/console.ts ***!
  \********************************/
/*! exports provided: default, LogLevel */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "LogLevel", function() { return LogLevel; });
/* harmony import */ var _random__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./random */ "./src/commons/random.ts");
/* harmony import */ var _collections_Dictionary__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../collections/Dictionary */ "./src/collections/Dictionary.ts");
var __spreadArrays = (undefined && undefined.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
/**
 * Extends standard console
 */


var LogLevel;
(function (LogLevel) {
    LogLevel[LogLevel["error"] = 0] = "error";
    LogLevel[LogLevel["warn"] = 1] = "warn";
    LogLevel[LogLevel["info"] = 2] = "info";
    LogLevel[LogLevel["debug"] = 3] = "debug";
})(LogLevel || (LogLevel = {}));
var console_ext = /** @class */ (function () {
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function console_ext() {
        this._uid = _random__WEBPACK_IMPORTED_MODULE_0__["default"].guid();
        this._level = LogLevel.info;
        this._class_levels = new _collections_Dictionary__WEBPACK_IMPORTED_MODULE_1__["Dictionary"]();
    }
    Object.defineProperty(console_ext.prototype, "uid", {
        // ------------------------------------------------------------------------
        //                      p r o p e r t i e s
        // ------------------------------------------------------------------------
        get: function () {
            return this._uid;
        },
        set: function (value) {
            this._uid = value;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(console_ext.prototype, "level", {
        get: function () {
            return this._level;
        },
        set: function (value) {
            this._level = value;
        },
        enumerable: true,
        configurable: true
    });
    // ------------------------------------------------------------------------
    //                      p u b l i c
    // ------------------------------------------------------------------------
    console_ext.prototype.setClassLevel = function (class_name, level) {
        this._class_levels.put(class_name, level);
    };
    console_ext.prototype.error = function (scope, error) {
        var args = [];
        for (var _i = 2; _i < arguments.length; _i++) {
            args[_i - 2] = arguments[_i];
        }
        console.error.apply(console, __spreadArrays(["[" + this.uid + "] " + scope, error], args));
    };
    ;
    console_ext.prototype.warn = function (scope) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        if (this.getLevel(scope) < LogLevel.warn) {
            return;
        }
        console.warn.apply(console, __spreadArrays(["[" + this.uid + "] " + scope], args));
    };
    ;
    console_ext.prototype.info = function (scope) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        if (this.getLevel(scope) < LogLevel.info) {
            return;
        }
        console.info.apply(console, __spreadArrays(["[" + this.uid + "] " + scope], args));
    };
    ;
    console_ext.prototype.debug = function (scope) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        if (this.getLevel(scope) < LogLevel.debug) {
            return;
        }
        console.log.apply(console, __spreadArrays(["[" + this.uid + "] " + scope], args));
    };
    ;
    console_ext.prototype.log = function (scope) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        if (this.getLevel(scope) < LogLevel.info) {
            return;
        }
        console.log.apply(console, __spreadArrays(["[" + this.uid + "] " + scope], args));
    };
    // ------------------------------------------------------------------------
    //                      p r i v a t e
    // ------------------------------------------------------------------------
    console_ext.prototype.getLevel = function (scope) {
        if (!!scope) {
            var class_name = scope.split(".")[0];
            if (this._class_levels.containsKey(class_name)) {
                return this._class_levels.get(class_name);
            }
        }
        return this.level;
    };
    console_ext.instance = function () {
        if (null == console_ext.__instance) {
            console_ext.__instance = new console_ext();
        }
        return console_ext.__instance;
    };
    return console_ext;
}());
// ------------------------------------------------------------------------
//                      e x p o r t
// ------------------------------------------------------------------------
/* harmony default export */ __webpack_exports__["default"] = (console_ext.instance());



/***/ }),

/***/ "./src/commons/lang.ts":
/*!*****************************!*\
  !*** ./src/commons/lang.ts ***!
  \*****************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/**
 * Utility class
 */
var __spreadArrays = (undefined && undefined.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
var langClass = /** @class */ (function () {
    // ------------------------------------------------------------------------
    //                      f i e l d s
    // ------------------------------------------------------------------------
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function langClass() {
    }
    Object.defineProperty(langClass.prototype, "window", {
        // ------------------------------------------------------------------------
        //                      p u b l i c
        // ------------------------------------------------------------------------
        get: function () {
            return window;
        },
        enumerable: true,
        configurable: true
    });
    langClass.prototype.parse = function (value) {
        try {
            if (this.isString(value)) {
                return JSON.parse(value);
            }
        }
        catch (err) {
        }
        return value;
    };
    // ------------------------------------------------------------------------
    //                      t o
    // ------------------------------------------------------------------------
    langClass.prototype.toString = function (value) {
        switch (typeof value) {
            case 'string':
            case 'number':
            case 'boolean':
                return value + '';
            case 'object':
                try {
                    // null is an object but is falsy. Swallow it.
                    return value === null ? '' : JSON.stringify(value);
                }
                catch (jsonError) {
                    return '{...}';
                }
            default:
                // Anything else will be replaced with an empty string
                // For example: undefined, Symbol, etc.
                return '';
        }
    };
    langClass.prototype.toArray = function (value) {
        return !!value
            ? this.isArray(value) ? value : [value]
            : [];
    };
    langClass.prototype.toBoolean = function (value, def_val) {
        return !!value
            ? value !== 'false' && value !== '0'
            : def_val;
    };
    langClass.prototype.toFloat = function (value, def_value, min, max) {
        if (def_value === void 0) { def_value = 0.0; }
        try {
            var result = parseFloat(value.replace(/,/g, '.'));
            result = this.isNaN(result) ? def_value : result;
            if (!this.isNaN(max) && result > (max || 0))
                result = max || 0;
            if (!this.isNaN(min) && result < (min || 0))
                result = min || 0;
            return result;
        }
        catch (err) {
            return def_value;
        }
    };
    langClass.prototype.toInt = function (value, def_value, min, max) {
        if (def_value === void 0) { def_value = 0; }
        try {
            var result = parseInt(value);
            result = this.isNaN(result) ? def_value : result;
            if (!this.isNaN(max) && result > (max || 0))
                result = max || 0;
            if (!this.isNaN(min) && result < (min || 0))
                result = min || 0;
            return result;
        }
        catch (err) {
            return def_value;
        }
    };
    // ------------------------------------------------------------------------
    //                      i s
    // ------------------------------------------------------------------------
    langClass.prototype.isFunction = function (value) {
        return typeof value == 'function';
    };
    langClass.prototype.isObject = function (value) {
        return value === Object(value);
    };
    langClass.prototype.isArray = function (value) {
        return !!Array.isArray
            ? Array.isArray(value)
            : value && typeof value == 'object' && typeof value.length == 'number' && toString.call(value) == '[object Array]' || false;
    };
    langClass.prototype.isArguments = function (value) {
        return value && typeof value == 'object' && typeof value.length == 'number' &&
            toString.call(value) == '[object Arguments]' || false;
    };
    langClass.prototype.isBoolean = function (value) {
        return value === true || value === false ||
            value && typeof value == 'object' && toString.call(value) == '[object Boolean]' || false;
    };
    langClass.prototype.isString = function (value) {
        return typeof value == 'string' ||
            value && typeof value == 'object' && toString.call(value) == '[object String]' || false;
    };
    langClass.prototype.isNumber = function (value) {
        return typeof value == 'number' ||
            value && typeof value == 'object' && toString.call(value) == '[object Number]' || false;
    };
    langClass.prototype.isNaN = function (value) {
        return isNaN(value);
    };
    langClass.isDate = function (value) {
        return value && typeof value == 'object' && toString.call(value) == '[object Date]' || false;
    };
    // ------------------------------------------------------------------------
    //                      u t i l s
    // ------------------------------------------------------------------------
    /**
     * Invoke a function. Shortcut for "func.call(this, ...args)"
     */
    langClass.prototype.funcInvoke = function (func) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        var self = this;
        if (!!func && self.isFunction(func)) {
            if (args.length === 0) {
                return func.call(self);
            }
            else {
                return func.call.apply(func, __spreadArrays([self], args));
            }
        }
        return null;
    };
    /**
     * Delays a function for the given number of milliseconds, and then calls
     * it with the arguments supplied.
     * NOTE: user "clearTimeout" with funcDelay response to
     */
    langClass.prototype.funcDelay = function (func, wait) {
        var args = [];
        for (var _i = 2; _i < arguments.length; _i++) {
            args[_i - 2] = arguments[_i];
        }
        return setTimeout(function () {
            return func.call.apply(func, __spreadArrays([null], args));
        }, wait);
    };
    langClass.instance = function () {
        if (null == langClass.__instance) {
            langClass.__instance = new langClass();
        }
        return langClass.__instance;
    };
    return langClass;
}());
// ------------------------------------------------------------------------
//                      e x p o r t
// ------------------------------------------------------------------------
var lang = langClass.instance();
/* harmony default export */ __webpack_exports__["default"] = (lang);


/***/ }),

/***/ "./src/commons/random.ts":
/*!*******************************!*\
  !*** ./src/commons/random.ts ***!
  \*******************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
var random = /** @class */ (function () {
    function random() {
    }
    // ------------------------------------------------------------------------
    //                      random and GUID
    // ------------------------------------------------------------------------
    // A (possibly faster) way to get the current timestamp as an integer.
    random.now = function () {
        return !!Date.now ? Date.now() : new Date().getTime();
    };
    random.guid = function () {
        return random._s4() + random._s4() + '-' + random._s4() + '-' + random._s4() + '-' +
            random._s4() + '-' + random._s4() + random._s4() + random._s4();
    };
    random.uniqueId = function (prefix) {
        var id = ++random._id_counter + '';
        return prefix ? prefix + id : id;
    };
    // ------------------------------------------------------------------------
    //                      p r i v a t e
    // ------------------------------------------------------------------------
    random._s4 = function () {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    };
    // ------------------------------------------------------------------------
    //                      f i e l d s
    // ------------------------------------------------------------------------
    random._id_counter = 0;
    return random;
}());
/* harmony default export */ __webpack_exports__["default"] = (random);


/***/ }),

/***/ "./src/events/EventEmitter.ts":
/*!************************************!*\
  !*** ./src/events/EventEmitter.ts ***!
  \************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _Events__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./Events */ "./src/events/Events.ts");
/* harmony import */ var _commons_BaseObject__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../commons/BaseObject */ "./src/commons/BaseObject.ts");
/* harmony import */ var _collections_Dictionary__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../collections/Dictionary */ "./src/collections/Dictionary.ts");
var __extends = (undefined && undefined.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
var __spreadArrays = (undefined && undefined.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};



/**
 * Class that emit events with a scope.
 */
var EventEmitter = /** @class */ (function (_super) {
    __extends(EventEmitter, _super);
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function EventEmitter() {
        var _this = _super.call(this) || this;
        _this._listeners = new _collections_Dictionary__WEBPACK_IMPORTED_MODULE_2__["Dictionary"]();
        return _this;
    }
    // ------------------------------------------------------------------------
    //                      p u b l i c
    // ------------------------------------------------------------------------
    EventEmitter.prototype.on = function (scope, eventName, listener) {
        if (!!scope) {
            var key = EventEmitter.key(scope);
            if (!this._listeners.containsKey(key)) {
                this._listeners.put(key, new _Events__WEBPACK_IMPORTED_MODULE_0__["default"]());
            }
            this._listeners.get(key).on(eventName, listener.bind(scope));
        }
    };
    EventEmitter.prototype.once = function (scope, eventName, listener) {
        if (!!scope) {
            var key = EventEmitter.key(scope);
            if (!this._listeners.containsKey(key)) {
                this._listeners.put(key, new _Events__WEBPACK_IMPORTED_MODULE_0__["default"]());
            }
            this._listeners.get(key).once(eventName, listener.bind(scope));
        }
    };
    EventEmitter.prototype.off = function (scope, eventName) {
        if (!!scope) {
            var key = EventEmitter.key(scope);
            if (this._listeners.containsKey(key)) {
                this._listeners.get(key).off(eventName);
            }
        }
    };
    EventEmitter.prototype.clear = function () {
        if (!!this._listeners) {
            var keys = this._listeners.keys();
            for (var _i = 0, keys_1 = keys; _i < keys_1.length; _i++) {
                var key = keys_1[_i];
                if (this._listeners.containsKey(key)) {
                    this._listeners.get(key).clear();
                }
            }
        }
    };
    EventEmitter.prototype.emit = function (eventName) {
        var _a;
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        if (!!this._listeners) {
            var keys = this._listeners.keys();
            for (var _b = 0, keys_2 = keys; _b < keys_2.length; _b++) {
                var key = keys_2[_b];
                if (this._listeners.containsKey(key)) {
                    (_a = this._listeners.get(key)).emit.apply(_a, __spreadArrays([eventName], args));
                }
            }
        }
    };
    // ------------------------------------------------------------------------
    //                      S T A T I C
    // ------------------------------------------------------------------------
    EventEmitter.key = function (scope) {
        try {
            return scope.uid;
        }
        catch (err) {
            console.warn("ApplicationEvents.key()", "BINDING EVENT ON DEFAULT KEY!");
            return '_default';
        }
    };
    return EventEmitter;
}(_commons_BaseObject__WEBPACK_IMPORTED_MODULE_1__["default"]));
/* harmony default export */ __webpack_exports__["default"] = (EventEmitter);


/***/ }),

/***/ "./src/events/Events.ts":
/*!******************************!*\
  !*** ./src/events/Events.ts ***!
  \******************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var _collections_Dictionary__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../collections/Dictionary */ "./src/collections/Dictionary.ts");
/* harmony import */ var _commons_lang__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../commons/lang */ "./src/commons/lang.ts");
/* harmony import */ var _commons_console__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../commons/console */ "./src/commons/console.ts");



/**
 * Events controller.
 *
 * <code>
 *
 * import {Events} from "./events/Events";
 *
 * class MyEmitter extends Events{}
 *
 * let myEmitter = new MyEmitter();
 * myEmitter.on('event', () => {
 *   console.log('event occured')
 * });
 *
 * myEmitter.emit('event');
 *
 * </code>
 *
 *
 */
var Events = /** @class */ (function () {
    function Events() {
        // ------------------------------------------------------------------------
        //                      C O N S T
        // ------------------------------------------------------------------------
        // ------------------------------------------------------------------------
        //                      f i e l d s
        // ------------------------------------------------------------------------
        this._events = new _collections_Dictionary__WEBPACK_IMPORTED_MODULE_0__["Dictionary"]();
        this._maxListeners = 0;
    }
    // ------------------------------------------------------------------------
    //                      p r o p e r t i e s
    // ------------------------------------------------------------------------
    Events.prototype.getMaxListeners = function () {
        return this._maxListeners === 0 ? Events.DEFAULT_MAX_LISTENERS : this._maxListeners;
    };
    Events.prototype.setMaxListeners = function (limit) {
        this._maxListeners = limit;
        return this;
    };
    // ------------------------------------------------------------------------
    //                      p u b l i c
    // ------------------------------------------------------------------------
    Events.prototype.addListener = function (eventName, listener) {
        return this.on(eventName, listener);
    };
    Events.prototype.on = function (eventName, listener) {
        this._registerEvent(eventName, listener, false);
        return this;
    };
    Events.prototype.once = function (eventName, listener) {
        this._registerEvent(eventName, listener, true);
        return this;
    };
    Events.prototype.off = function (event_names, listener) {
        var names = _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].isArray(event_names)
            ? event_names
            : !!event_names ? [event_names] : [];
        if (!!listener) {
            for (var _i = 0, names_1 = names; _i < names_1.length; _i++) {
                var name_1 = names_1[_i];
                this.removeListener(name_1, listener);
            }
        }
        else {
            if (names.length > 0) {
                this.removeAllListeners(names);
            }
            else {
                this.removeAllListeners();
            }
        }
        return this;
    };
    Events.prototype.emit = function (eventName) {
        var args = [];
        for (var _i = 1; _i < arguments.length; _i++) {
            args[_i - 1] = arguments[_i];
        }
        var listeners = this._events.get(eventName);
        var listenerCount = this.listenerCount(eventName);
        if (listeners) {
            listeners.map(function (listener) { return listener.apply(void 0, args); });
        }
        return listenerCount !== 0;
    };
    Events.prototype.eventNames = function () {
        return this._events.keys();
    };
    Events.prototype.listeners = function (eventName) {
        return this._events.get(eventName);
    };
    Events.prototype.listenerCount = function (eventName) {
        var listeners = this._events.get(eventName);
        return listeners === undefined ? 0 : listeners.length;
    };
    Events.prototype.removeAllListeners = function (eventNames) {
        var _this = this;
        if (!eventNames) {
            eventNames = this._events.keys();
        }
        eventNames.forEach(function (eventName) { return _this._events.remove(eventName); });
        return this;
    };
    Events.prototype.removeListener = function (eventName, listener) {
        var listeners = this.listeners(eventName);
        var filtered_listeners = !!listeners
            ? listeners.filter(function (item) { return item === listener; }) // filter only valid
            : [];
        this._events.put(eventName, filtered_listeners);
        return this;
    };
    Events.prototype.clear = function () {
        this._events.clear();
    };
    // ------------------------------------------------------------------------
    //                      p r i v a t e
    // ------------------------------------------------------------------------
    Events.prototype._registerEvent = function (eventName, listener, type) {
        if (this._listenerLimitReached(eventName)) {
            _commons_console__WEBPACK_IMPORTED_MODULE_2__["default"].warn("Events._registerEvent", "Maximum listener reached, new Listener not added", this.getMaxListeners());
            return;
        }
        if (type === true) {
            listener = this._createOnceListener(listener, eventName);
        }
        var listeners = Events._createListeners(listener, this.listeners(eventName));
        this._events.put(eventName, listeners);
        return;
    };
    Events.prototype._createOnceListener = function (listener, eventName) {
        var _this = this;
        return function () {
            var args = [];
            for (var _i = 0; _i < arguments.length; _i++) {
                args[_i] = arguments[_i];
            }
            _this.removeListener(eventName, listener);
            return listener.apply(void 0, args);
        };
    };
    Events.prototype._listenerLimitReached = function (eventName) {
        return this.listenerCount(eventName) >= this.getMaxListeners();
    };
    Events._createListeners = function (listener, listeners) {
        if (!listeners) {
            listeners = [];
        }
        listeners.push(listener);
        return listeners;
    };
    Events.DEFAULT_MAX_LISTENERS = 10; // max listener for each event name
    return Events;
}());
// ------------------------------------------------------------------------
//                      e x p o r t s
// ------------------------------------------------------------------------
/* harmony default export */ __webpack_exports__["default"] = (Events);


/***/ }),

/***/ "./src/launcher.ts":
/*!*************************!*\
  !*** ./src/launcher.ts ***!
  \*************************/
/*! exports provided: ROOT */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ROOT", function() { return ROOT; });
/* harmony import */ var _commons_BaseObject__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./commons/BaseObject */ "./src/commons/BaseObject.ts");
/* harmony import */ var _commons_console__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./commons/console */ "./src/commons/console.ts");
/* harmony import */ var _socket_SocketService__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./socket/SocketService */ "./src/socket/SocketService.ts");
/* harmony import */ var _commons_lang__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./commons/lang */ "./src/commons/lang.ts");
var __extends = (undefined && undefined.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();




// ------------------------------------------------------------------------
//                      c o n s t a n t s
// ------------------------------------------------------------------------
var DEBUG = false;
var VERSION = "1.0.1";
var UID = "__vws";
var ROOT = window;
/**
 * Launcher class
 */
var launcher = /** @class */ (function (_super) {
    __extends(launcher, _super);
    // ------------------------------------------------------------------------
    //                      f i e l d s
    // ------------------------------------------------------------------------
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function launcher() {
        var _this = _super.call(this) || this;
        _this.init();
        return _this;
    }
    // ------------------------------------------------------------------------
    //                      p u b l i c
    // ------------------------------------------------------------------------
    launcher.prototype.start = function () {
        _commons_console__WEBPACK_IMPORTED_MODULE_1__["default"].debug("launcher.start()");
        // creates Vanilla.Websocket builder reference
        ROOT[UID] = {
            version: VERSION,
            create: this.newClient,
        };
    };
    // ------------------------------------------------------------------------
    //                      p r i v a t e
    // ------------------------------------------------------------------------
    launcher.prototype.init = function () {
        // init application scope
        _commons_BaseObject__WEBPACK_IMPORTED_MODULE_0__["default"].PREFIX = UID + "_"; // application uid become component prefix.
        // set console prefix
        _commons_console__WEBPACK_IMPORTED_MODULE_1__["default"].uid = UID;
        if (DEBUG) {
            _commons_console__WEBPACK_IMPORTED_MODULE_1__["default"].level = _commons_console__WEBPACK_IMPORTED_MODULE_1__["LogLevel"].debug;
        }
        else {
            _commons_console__WEBPACK_IMPORTED_MODULE_1__["default"].level = _commons_console__WEBPACK_IMPORTED_MODULE_1__["LogLevel"].info;
        }
    };
    launcher.prototype.newClient = function (params) {
        if (_commons_lang__WEBPACK_IMPORTED_MODULE_3__["default"].isString(params)) {
            params = {
                host: params
            };
        }
        else {
            params = params || {};
        }
        return new _socket_SocketService__WEBPACK_IMPORTED_MODULE_2__["SocketService"](params);
    };
    launcher.instance = function () {
        if (null == launcher.__instance) {
            launcher.__instance = new launcher();
        }
        return launcher.__instance;
    };
    return launcher;
}(_commons_BaseObject__WEBPACK_IMPORTED_MODULE_0__["default"]));
// ------------------------------------------------------------------------
//                      S T A R T   A P P L I C A T I O N
// ------------------------------------------------------------------------
launcher.instance().start();
// ------------------------------------------------------------------------
//                      e x p o r t
// ------------------------------------------------------------------------



/***/ }),

/***/ "./src/socket/SocketService.ts":
/*!*************************************!*\
  !*** ./src/socket/SocketService.ts ***!
  \*************************************/
/*! exports provided: SocketService, EVENT_CLOSE, EVENT_MESSAGE, EVENT_OPEN, EVENT_ERROR */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "SocketService", function() { return SocketService; });
/* harmony import */ var _events_EventEmitter__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../events/EventEmitter */ "./src/events/EventEmitter.ts");
/* harmony import */ var _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./WebSocketChannel */ "./src/socket/WebSocketChannel.ts");
/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "EVENT_CLOSE", function() { return _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_CLOSE"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "EVENT_MESSAGE", function() { return _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_MESSAGE"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "EVENT_OPEN", function() { return _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_OPEN"]; });

/* harmony reexport (safe) */ __webpack_require__.d(__webpack_exports__, "EVENT_ERROR", function() { return _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_ERROR"]; });

/* harmony import */ var _commons_lang__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../commons/lang */ "./src/commons/lang.ts");
var __extends = (undefined && undefined.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();



var SocketService = /** @class */ (function (_super) {
    __extends(SocketService, _super);
    // ------------------------------------------------------------------------
    //  c o n s t r u c t o r
    // ------------------------------------------------------------------------
    function SocketService(params) {
        var _this = _super.call(this) || this;
        if (_commons_lang__WEBPACK_IMPORTED_MODULE_2__["default"].isString(params)) {
            _this._params = {
                host: params
            };
        }
        else {
            _this._params = params || {};
        }
        _this._is_connected = false;
        return _this;
    }
    Object.defineProperty(SocketService.prototype, "host", {
        // ------------------------------------------------------------------------
        //  p u b l i c
        // ------------------------------------------------------------------------
        get: function () {
            return !!this._params ? this._params["host"] || "" : "";
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(SocketService.prototype, "isConnected", {
        get: function () {
            return this._is_connected;
        },
        enumerable: true,
        configurable: true
    });
    SocketService.prototype.send = function (message, callback) {
        var _this = this;
        this.socket.open().ready(function (is_ready, error) {
            if (is_ready) {
                _this.socket.send(message, callback);
            }
            else {
                _commons_lang__WEBPACK_IMPORTED_MODULE_2__["default"].funcInvoke(callback, { "error": error });
            }
        });
    };
    SocketService.prototype.close = function () {
        if (!!this._socket) {
            if (this.socket.initialized) {
                this.socket.close();
            }
        }
    };
    SocketService.prototype.reset = function () {
        if (!!this._socket) {
            if (this.socket.initialized) {
                this.socket.reset();
            }
        }
    };
    // @ts-ignore
    SocketService.prototype.on = function (eventName, listener) {
        _super.prototype.on.call(this, this, eventName, listener);
    };
    // @ts-ignore
    SocketService.prototype.off = function (eventName) {
        _super.prototype.off.call(this, this, eventName);
    };
    Object.defineProperty(SocketService.prototype, "socket", {
        // ------------------------------------------------------------------------
        //  p r i v a t e
        // ------------------------------------------------------------------------
        get: function () {
            if (!this._socket) {
                // creates new socket
                this._socket = new _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["WebSocketChannel"](this._params);
                this._socket.on(this, _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_OPEN"], this.onOpen);
                this._socket.on(this, _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_CLOSE"], this.onClose);
                this._socket.on(this, _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_MESSAGE"], this.onMessage);
                this._socket.on(this, _WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_ERROR"], this.onError);
            }
            return this._socket;
        },
        enumerable: true,
        configurable: true
    });
    SocketService.prototype.onOpen = function () {
        this._is_connected = true;
        this.emit(_WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_OPEN"]);
    };
    SocketService.prototype.onClose = function () {
        this._is_connected = false;
        this.emit(_WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_CLOSE"]);
    };
    SocketService.prototype.onMessage = function (data) {
        this.emit(_WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_MESSAGE"], data);
    };
    SocketService.prototype.onError = function (error) {
        this._is_connected = false;
        this.emit(_WebSocketChannel__WEBPACK_IMPORTED_MODULE_1__["EVENT_ERROR"], error);
    };
    return SocketService;
}(_events_EventEmitter__WEBPACK_IMPORTED_MODULE_0__["default"]));
// ------------------------------------------------------------------------
//  E X P O R T
// ------------------------------------------------------------------------



/***/ }),

/***/ "./src/socket/WebSocketChannel.ts":
/*!****************************************!*\
  !*** ./src/socket/WebSocketChannel.ts ***!
  \****************************************/
/*! exports provided: WebSocketChannel, EVENT_CLOSE, EVENT_MESSAGE, EVENT_OPEN, EVENT_ERROR */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "WebSocketChannel", function() { return WebSocketChannel; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "EVENT_CLOSE", function() { return EVENT_CLOSE; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "EVENT_MESSAGE", function() { return EVENT_MESSAGE; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "EVENT_OPEN", function() { return EVENT_OPEN; });
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "EVENT_ERROR", function() { return EVENT_ERROR; });
/* harmony import */ var _events_EventEmitter__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ../events/EventEmitter */ "./src/events/EventEmitter.ts");
/* harmony import */ var _commons_lang__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ../commons/lang */ "./src/commons/lang.ts");
/* harmony import */ var _commons_random__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ../commons/random */ "./src/commons/random.ts");
var __extends = (undefined && undefined.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();



var DEF_HOST = 'ws://localhost:8181/websocket';
var DEF_UID = 'mebot_web';
var EVENT_OPEN = 'on_open';
var EVENT_CLOSE = 'on_close';
var EVENT_MESSAGE = 'on_message';
var EVENT_ERROR = 'on_error';
var FLD_REQUEST_UUID = "request_uuid";
var FLD_REQUEST_UUID_HANDLED = "request_uuid_handled";
var FLD_REQUEST_UUID_TIMEOUT = 10 * 1000; // 10 seconds timeout
var WebSocketChannel = /** @class */ (function (_super) {
    __extends(WebSocketChannel, _super);
    // ------------------------------------------------------------------------
    //                      c o n s t r u c t o r
    // ------------------------------------------------------------------------
    /**
     * Creates a WebSocket wrapper
     * @param params "{host:'ws://localhost:8181/websocket'}"
     */
    function WebSocketChannel(params) {
        var _this = _super.call(this) || this;
        _this._initialized = false;
        _this._active = false;
        _this._host = !!params ? params.host || DEF_HOST : DEF_HOST;
        _this._callback_pool = {}; // contains registered callbacks
        return _this;
    }
    // ------------------------------------------------------------------------
    //                      p u b l i c
    // ------------------------------------------------------------------------
    /**
     * Use this to pass a callback and enable message send when
     * socket is ready to send.
     */
    WebSocketChannel.prototype.ready = function (callback) {
        var _this = this;
        if (this._initialized) {
            if (this._active) {
                _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].funcInvoke(callback, true);
            }
            else {
                // timeout for ready status  (3 seconds)
                _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].funcDelay(function () {
                    if (!_this._active) {
                        _this.off(_this); // clear buffer
                        _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].funcInvoke(callback, false, "timeout");
                        // reset socket status
                        _this.free();
                    }
                }, 3 * 1000);
                this.off(this); // clear buffer
                this.on(this, EVENT_OPEN, function () {
                    _this.off(_this); // clear buffer
                    _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].funcInvoke(callback, true);
                });
            }
        }
        else {
            // exit not ready
            _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].funcInvoke(callback, false, "not initialized");
        }
    };
    WebSocketChannel.prototype.reset = function () {
        this._callback_pool = {};
        this.close();
        this.open();
    };
    Object.defineProperty(WebSocketChannel.prototype, "initialized", {
        /**
         * Socket is properly configured
         */
        get: function () {
            return this._initialized;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(WebSocketChannel.prototype, "active", {
        /**
         * Socket is Open
         */
        get: function () {
            return this._active;
        },
        enumerable: true,
        configurable: true
    });
    Object.defineProperty(WebSocketChannel.prototype, "host", {
        get: function () {
            return this._host;
        },
        enumerable: true,
        configurable: true
    });
    WebSocketChannel.prototype.open = function () {
        if (!this._web_socket) {
            try {
                this._web_socket = this.createWs();
                if (this.handle(this._web_socket)) {
                    this._initialized = true;
                }
                else {
                    this._initialized = false;
                }
            }
            catch (err) {
                console.error('WebSocketChannel.open', err);
            }
        }
        return this;
    };
    WebSocketChannel.prototype.close = function () {
        try {
            if (this._active) {
                this._initialized = false;
                this._active = false;
                // close and free socket
                this.free();
            }
        }
        catch (err) {
            console.error('WebSocketChannel.close', err);
        }
    };
    WebSocketChannel.prototype.send = function (message, callback) {
        var _this = this;
        try {
            if (this._active
                && !!this._web_socket
                && this._web_socket.readyState === this._web_socket.OPEN) {
                if (!!callback) {
                    var callback_uuid_1 = _commons_random__WEBPACK_IMPORTED_MODULE_2__["default"].guid();
                    message[FLD_REQUEST_UUID] = callback_uuid_1;
                    this._callback_pool[callback_uuid_1] = callback;
                    // set timeout
                    _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].funcDelay(function () {
                        delete _this._callback_pool[callback_uuid_1];
                        console.debug("WebSocketChannel.send", "timeout removed " + callback_uuid_1);
                    }, FLD_REQUEST_UUID_TIMEOUT);
                }
                // ready to send
                this._web_socket.send(_commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].toString(message));
            }
            else {
                console.warn('WebSocketChannel.send', 'Socket is not ready to send message.', this._web_socket, message);
            }
        }
        catch (err) {
            console.error('WebSocketChannel.send', err);
        }
    };
    // ------------------------------------------------------------------------
    //                      p r i v a t e
    // ------------------------------------------------------------------------
    WebSocketChannel.prototype.createWs = function () {
        var WS_native = _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].window['WebSocket'] || _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].window['MozWebSocket'];
        return new WS_native(this._host);
    };
    WebSocketChannel.prototype.free = function () {
        if (!!this._web_socket) {
            try {
                this._web_socket.close();
            }
            catch (err) {
                console.debug("WebSocketChannel.free", err);
            }
            this._web_socket = null;
        }
    };
    WebSocketChannel.prototype.handle = function (ws) {
        if (!!ws) {
            ws.onmessage = this._on_message.bind(this);
            ws.onopen = this._on_open.bind(this);
            ws.onclose = this._on_close.bind(this);
            ws.onerror = this._on_error.bind(this);
            return true;
        }
        else {
            console.warn('WebSocketChannel.handle', 'WebSocket not found', ws);
            return false;
        }
    };
    WebSocketChannel.prototype._on_open = function (ev) {
        this._active = true;
        this.emit(EVENT_OPEN);
    };
    WebSocketChannel.prototype._on_close = function (ev) {
        this._active = false;
        this.emit(EVENT_CLOSE);
    };
    WebSocketChannel.prototype._on_message = function (ev) {
        try {
            var origin_1 = ev.origin;
            var ports = ev.ports;
            var data = _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].parse(ev.data);
            if (!!data[FLD_REQUEST_UUID]) {
                var callback_uuid = data[FLD_REQUEST_UUID];
                if (!!this._callback_pool[callback_uuid]) {
                    var f = this._callback_pool[callback_uuid];
                    if (_commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].isFunction(f)) {
                        f(data);
                    }
                }
                // remove
                delete this._callback_pool[callback_uuid];
                data[FLD_REQUEST_UUID_HANDLED] = true;
            }
            this.emit(EVENT_MESSAGE, data);
        }
        catch (err) {
            console.error("WebSocketChannel._on_message", err);
        }
    };
    WebSocketChannel.prototype._on_error = function (ev) {
        ev.preventDefault();
        var str_err = _commons_lang__WEBPACK_IMPORTED_MODULE_1__["default"].toString(ev);
        console.error('WebSocketChannel._on_error', str_err);
        this.emit(EVENT_ERROR, str_err);
    };
    return WebSocketChannel;
}(_events_EventEmitter__WEBPACK_IMPORTED_MODULE_0__["default"]));



/***/ })

/******/ });
//# sourceMappingURL=vws.js.map