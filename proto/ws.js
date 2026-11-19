/*eslint-disable block-scoped-var, id-length, no-control-regex, no-magic-numbers, no-prototype-builtins, no-redeclare, no-shadow, no-var, sort-vars*/
import * as $protobuf from "protobufjs/minimal";

// Common aliases
const $Reader = $protobuf.Reader, $Writer = $protobuf.Writer, $util = $protobuf.util;

// Exported root namespace
const $root = $protobuf.roots["default"] || ($protobuf.roots["default"] = {});

export const ws = $root.ws = (() => {

    /**
     * Namespace ws.
     * @exports ws
     * @namespace
     */
    const ws = {};

    ws.WSReq = (function() {

        /**
         * Properties of a WSReq.
         * @memberof ws
         * @interface IWSReq
         * @property {number|Long|null} [seq] WSReq seq
         * @property {string|null} [name] WSReq name
         * @property {string|null} [proxy] WSReq proxy
         * @property {string|null} [authUid] WSReq authUid
         * @property {string|null} [authUsername] WSReq authUsername
         * @property {ws.IAccountLogin|null} [accountLogin] WSReq accountLogin
         * @property {ws.IAccountRegister|null} [accountRegister] WSReq accountRegister
         * @property {ws.IChatWorldSend|null} [chatWorldSend] WSReq chatWorldSend
         * @property {ws.IChatPrivateSend|null} [chatPrivateSend] WSReq chatPrivateSend
         * @property {ws.IChatGroupSend|null} [chatGroupSend] WSReq chatGroupSend
         * @property {ws.IGroupMember|null} [groupMemberReq] WSReq groupMemberReq
         * @property {ws.IRoomCreate|null} [roomCreate] WSReq roomCreate
         * @property {ws.IRoomJoin|null} [roomJoin] WSReq roomJoin
         * @property {ws.IRoomLeave|null} [roomLeave] WSReq roomLeave
         * @property {ws.IRoomKick|null} [roomKick] WSReq roomKick
         * @property {ws.IRoomDelete|null} [roomDelete] WSReq roomDelete
         * @property {ws.IRoomSend|null} [roomSend] WSReq roomSend
         * @property {ws.IInventoryAdd|null} [inventoryAdd] WSReq inventoryAdd
         * @property {ws.IInventoryConsume|null} [inventoryConsume] WSReq inventoryConsume
         * @property {ws.IInventoryDrop|null} [inventoryDrop] WSReq inventoryDrop
         * @property {ws.IInventoryGet|null} [inventoryGet] WSReq inventoryGet
         * @property {string|null} [msgJson] WSReq msgJson
         */

        /**
         * Constructs a new WSReq.
         * @memberof ws
         * @classdesc Represents a WSReq.
         * @implements IWSReq
         * @constructor
         * @param {ws.IWSReq=} [properties] Properties to set
         */
        function WSReq(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * WSReq seq.
         * @member {number|Long} seq
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.seq = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * WSReq name.
         * @member {string} name
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.name = "";

        /**
         * WSReq proxy.
         * @member {string} proxy
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.proxy = "";

        /**
         * WSReq authUid.
         * @member {string} authUid
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.authUid = "";

        /**
         * WSReq authUsername.
         * @member {string} authUsername
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.authUsername = "";

        /**
         * WSReq accountLogin.
         * @member {ws.IAccountLogin|null|undefined} accountLogin
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.accountLogin = null;

        /**
         * WSReq accountRegister.
         * @member {ws.IAccountRegister|null|undefined} accountRegister
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.accountRegister = null;

        /**
         * WSReq chatWorldSend.
         * @member {ws.IChatWorldSend|null|undefined} chatWorldSend
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.chatWorldSend = null;

        /**
         * WSReq chatPrivateSend.
         * @member {ws.IChatPrivateSend|null|undefined} chatPrivateSend
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.chatPrivateSend = null;

        /**
         * WSReq chatGroupSend.
         * @member {ws.IChatGroupSend|null|undefined} chatGroupSend
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.chatGroupSend = null;

        /**
         * WSReq groupMemberReq.
         * @member {ws.IGroupMember|null|undefined} groupMemberReq
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.groupMemberReq = null;

        /**
         * WSReq roomCreate.
         * @member {ws.IRoomCreate|null|undefined} roomCreate
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.roomCreate = null;

        /**
         * WSReq roomJoin.
         * @member {ws.IRoomJoin|null|undefined} roomJoin
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.roomJoin = null;

        /**
         * WSReq roomLeave.
         * @member {ws.IRoomLeave|null|undefined} roomLeave
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.roomLeave = null;

        /**
         * WSReq roomKick.
         * @member {ws.IRoomKick|null|undefined} roomKick
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.roomKick = null;

        /**
         * WSReq roomDelete.
         * @member {ws.IRoomDelete|null|undefined} roomDelete
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.roomDelete = null;

        /**
         * WSReq roomSend.
         * @member {ws.IRoomSend|null|undefined} roomSend
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.roomSend = null;

        /**
         * WSReq inventoryAdd.
         * @member {ws.IInventoryAdd|null|undefined} inventoryAdd
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.inventoryAdd = null;

        /**
         * WSReq inventoryConsume.
         * @member {ws.IInventoryConsume|null|undefined} inventoryConsume
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.inventoryConsume = null;

        /**
         * WSReq inventoryDrop.
         * @member {ws.IInventoryDrop|null|undefined} inventoryDrop
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.inventoryDrop = null;

        /**
         * WSReq inventoryGet.
         * @member {ws.IInventoryGet|null|undefined} inventoryGet
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.inventoryGet = null;

        /**
         * WSReq msgJson.
         * @member {string} msgJson
         * @memberof ws.WSReq
         * @instance
         */
        WSReq.prototype.msgJson = "";

        // OneOf field names bound to virtual getters and setters
        let $oneOfFields;

        /**
         * WSReq msg.
         * @member {"accountLogin"|"accountRegister"|"chatWorldSend"|"chatPrivateSend"|"chatGroupSend"|"groupMemberReq"|"roomCreate"|"roomJoin"|"roomLeave"|"roomKick"|"roomDelete"|"roomSend"|"inventoryAdd"|"inventoryConsume"|"inventoryDrop"|"inventoryGet"|undefined} msg
         * @memberof ws.WSReq
         * @instance
         */
        Object.defineProperty(WSReq.prototype, "msg", {
            get: $util.oneOfGetter($oneOfFields = ["accountLogin", "accountRegister", "chatWorldSend", "chatPrivateSend", "chatGroupSend", "groupMemberReq", "roomCreate", "roomJoin", "roomLeave", "roomKick", "roomDelete", "roomSend", "inventoryAdd", "inventoryConsume", "inventoryDrop", "inventoryGet"]),
            set: $util.oneOfSetter($oneOfFields)
        });

        /**
         * Creates a new WSReq instance using the specified properties.
         * @function create
         * @memberof ws.WSReq
         * @static
         * @param {ws.IWSReq=} [properties] Properties to set
         * @returns {ws.WSReq} WSReq instance
         */
        WSReq.create = function create(properties) {
            return new WSReq(properties);
        };

        /**
         * Encodes the specified WSReq message. Does not implicitly {@link ws.WSReq.verify|verify} messages.
         * @function encode
         * @memberof ws.WSReq
         * @static
         * @param {ws.IWSReq} message WSReq message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        WSReq.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.seq != null && Object.hasOwnProperty.call(message, "seq"))
                writer.uint32(/* id 1, wireType 0 =*/8).int64(message.seq);
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
            if (message.proxy != null && Object.hasOwnProperty.call(message, "proxy"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.proxy);
            if (message.authUid != null && Object.hasOwnProperty.call(message, "authUid"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.authUid);
            if (message.authUsername != null && Object.hasOwnProperty.call(message, "authUsername"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.authUsername);
            if (message.accountLogin != null && Object.hasOwnProperty.call(message, "accountLogin"))
                $root.ws.AccountLogin.encode(message.accountLogin, writer.uint32(/* id 10, wireType 2 =*/82).fork()).ldelim();
            if (message.accountRegister != null && Object.hasOwnProperty.call(message, "accountRegister"))
                $root.ws.AccountRegister.encode(message.accountRegister, writer.uint32(/* id 11, wireType 2 =*/90).fork()).ldelim();
            if (message.chatWorldSend != null && Object.hasOwnProperty.call(message, "chatWorldSend"))
                $root.ws.ChatWorldSend.encode(message.chatWorldSend, writer.uint32(/* id 20, wireType 2 =*/162).fork()).ldelim();
            if (message.chatPrivateSend != null && Object.hasOwnProperty.call(message, "chatPrivateSend"))
                $root.ws.ChatPrivateSend.encode(message.chatPrivateSend, writer.uint32(/* id 21, wireType 2 =*/170).fork()).ldelim();
            if (message.chatGroupSend != null && Object.hasOwnProperty.call(message, "chatGroupSend"))
                $root.ws.ChatGroupSend.encode(message.chatGroupSend, writer.uint32(/* id 22, wireType 2 =*/178).fork()).ldelim();
            if (message.groupMemberReq != null && Object.hasOwnProperty.call(message, "groupMemberReq"))
                $root.ws.GroupMember.encode(message.groupMemberReq, writer.uint32(/* id 23, wireType 2 =*/186).fork()).ldelim();
            if (message.roomCreate != null && Object.hasOwnProperty.call(message, "roomCreate"))
                $root.ws.RoomCreate.encode(message.roomCreate, writer.uint32(/* id 30, wireType 2 =*/242).fork()).ldelim();
            if (message.roomJoin != null && Object.hasOwnProperty.call(message, "roomJoin"))
                $root.ws.RoomJoin.encode(message.roomJoin, writer.uint32(/* id 31, wireType 2 =*/250).fork()).ldelim();
            if (message.roomLeave != null && Object.hasOwnProperty.call(message, "roomLeave"))
                $root.ws.RoomLeave.encode(message.roomLeave, writer.uint32(/* id 32, wireType 2 =*/258).fork()).ldelim();
            if (message.roomKick != null && Object.hasOwnProperty.call(message, "roomKick"))
                $root.ws.RoomKick.encode(message.roomKick, writer.uint32(/* id 33, wireType 2 =*/266).fork()).ldelim();
            if (message.roomDelete != null && Object.hasOwnProperty.call(message, "roomDelete"))
                $root.ws.RoomDelete.encode(message.roomDelete, writer.uint32(/* id 34, wireType 2 =*/274).fork()).ldelim();
            if (message.roomSend != null && Object.hasOwnProperty.call(message, "roomSend"))
                $root.ws.RoomSend.encode(message.roomSend, writer.uint32(/* id 35, wireType 2 =*/282).fork()).ldelim();
            if (message.inventoryAdd != null && Object.hasOwnProperty.call(message, "inventoryAdd"))
                $root.ws.InventoryAdd.encode(message.inventoryAdd, writer.uint32(/* id 50, wireType 2 =*/402).fork()).ldelim();
            if (message.inventoryConsume != null && Object.hasOwnProperty.call(message, "inventoryConsume"))
                $root.ws.InventoryConsume.encode(message.inventoryConsume, writer.uint32(/* id 51, wireType 2 =*/410).fork()).ldelim();
            if (message.inventoryDrop != null && Object.hasOwnProperty.call(message, "inventoryDrop"))
                $root.ws.InventoryDrop.encode(message.inventoryDrop, writer.uint32(/* id 52, wireType 2 =*/418).fork()).ldelim();
            if (message.inventoryGet != null && Object.hasOwnProperty.call(message, "inventoryGet"))
                $root.ws.InventoryGet.encode(message.inventoryGet, writer.uint32(/* id 53, wireType 2 =*/426).fork()).ldelim();
            if (message.msgJson != null && Object.hasOwnProperty.call(message, "msgJson"))
                writer.uint32(/* id 100, wireType 2 =*/802).string(message.msgJson);
            return writer;
        };

        /**
         * Encodes the specified WSReq message, length delimited. Does not implicitly {@link ws.WSReq.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.WSReq
         * @static
         * @param {ws.IWSReq} message WSReq message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        WSReq.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a WSReq message from the specified reader or buffer.
         * @function decode
         * @memberof ws.WSReq
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.WSReq} WSReq
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        WSReq.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.WSReq();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.seq = reader.int64();
                    break;
                case 2:
                    message.name = reader.string();
                    break;
                case 3:
                    message.proxy = reader.string();
                    break;
                case 4:
                    message.authUid = reader.string();
                    break;
                case 5:
                    message.authUsername = reader.string();
                    break;
                case 10:
                    message.accountLogin = $root.ws.AccountLogin.decode(reader, reader.uint32());
                    break;
                case 11:
                    message.accountRegister = $root.ws.AccountRegister.decode(reader, reader.uint32());
                    break;
                case 20:
                    message.chatWorldSend = $root.ws.ChatWorldSend.decode(reader, reader.uint32());
                    break;
                case 21:
                    message.chatPrivateSend = $root.ws.ChatPrivateSend.decode(reader, reader.uint32());
                    break;
                case 22:
                    message.chatGroupSend = $root.ws.ChatGroupSend.decode(reader, reader.uint32());
                    break;
                case 23:
                    message.groupMemberReq = $root.ws.GroupMember.decode(reader, reader.uint32());
                    break;
                case 30:
                    message.roomCreate = $root.ws.RoomCreate.decode(reader, reader.uint32());
                    break;
                case 31:
                    message.roomJoin = $root.ws.RoomJoin.decode(reader, reader.uint32());
                    break;
                case 32:
                    message.roomLeave = $root.ws.RoomLeave.decode(reader, reader.uint32());
                    break;
                case 33:
                    message.roomKick = $root.ws.RoomKick.decode(reader, reader.uint32());
                    break;
                case 34:
                    message.roomDelete = $root.ws.RoomDelete.decode(reader, reader.uint32());
                    break;
                case 35:
                    message.roomSend = $root.ws.RoomSend.decode(reader, reader.uint32());
                    break;
                case 50:
                    message.inventoryAdd = $root.ws.InventoryAdd.decode(reader, reader.uint32());
                    break;
                case 51:
                    message.inventoryConsume = $root.ws.InventoryConsume.decode(reader, reader.uint32());
                    break;
                case 52:
                    message.inventoryDrop = $root.ws.InventoryDrop.decode(reader, reader.uint32());
                    break;
                case 53:
                    message.inventoryGet = $root.ws.InventoryGet.decode(reader, reader.uint32());
                    break;
                case 100:
                    message.msgJson = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a WSReq message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.WSReq
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.WSReq} WSReq
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        WSReq.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a WSReq message.
         * @function verify
         * @memberof ws.WSReq
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        WSReq.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            let properties = {};
            if (message.seq != null && message.hasOwnProperty("seq"))
                if (!$util.isInteger(message.seq) && !(message.seq && $util.isInteger(message.seq.low) && $util.isInteger(message.seq.high)))
                    return "seq: integer|Long expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            if (message.proxy != null && message.hasOwnProperty("proxy"))
                if (!$util.isString(message.proxy))
                    return "proxy: string expected";
            if (message.authUid != null && message.hasOwnProperty("authUid"))
                if (!$util.isString(message.authUid))
                    return "authUid: string expected";
            if (message.authUsername != null && message.hasOwnProperty("authUsername"))
                if (!$util.isString(message.authUsername))
                    return "authUsername: string expected";
            if (message.accountLogin != null && message.hasOwnProperty("accountLogin")) {
                properties.msg = 1;
                {
                    let error = $root.ws.AccountLogin.verify(message.accountLogin);
                    if (error)
                        return "accountLogin." + error;
                }
            }
            if (message.accountRegister != null && message.hasOwnProperty("accountRegister")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.AccountRegister.verify(message.accountRegister);
                    if (error)
                        return "accountRegister." + error;
                }
            }
            if (message.chatWorldSend != null && message.hasOwnProperty("chatWorldSend")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatWorldSend.verify(message.chatWorldSend);
                    if (error)
                        return "chatWorldSend." + error;
                }
            }
            if (message.chatPrivateSend != null && message.hasOwnProperty("chatPrivateSend")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatPrivateSend.verify(message.chatPrivateSend);
                    if (error)
                        return "chatPrivateSend." + error;
                }
            }
            if (message.chatGroupSend != null && message.hasOwnProperty("chatGroupSend")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatGroupSend.verify(message.chatGroupSend);
                    if (error)
                        return "chatGroupSend." + error;
                }
            }
            if (message.groupMemberReq != null && message.hasOwnProperty("groupMemberReq")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.GroupMember.verify(message.groupMemberReq);
                    if (error)
                        return "groupMemberReq." + error;
                }
            }
            if (message.roomCreate != null && message.hasOwnProperty("roomCreate")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomCreate.verify(message.roomCreate);
                    if (error)
                        return "roomCreate." + error;
                }
            }
            if (message.roomJoin != null && message.hasOwnProperty("roomJoin")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomJoin.verify(message.roomJoin);
                    if (error)
                        return "roomJoin." + error;
                }
            }
            if (message.roomLeave != null && message.hasOwnProperty("roomLeave")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomLeave.verify(message.roomLeave);
                    if (error)
                        return "roomLeave." + error;
                }
            }
            if (message.roomKick != null && message.hasOwnProperty("roomKick")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomKick.verify(message.roomKick);
                    if (error)
                        return "roomKick." + error;
                }
            }
            if (message.roomDelete != null && message.hasOwnProperty("roomDelete")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomDelete.verify(message.roomDelete);
                    if (error)
                        return "roomDelete." + error;
                }
            }
            if (message.roomSend != null && message.hasOwnProperty("roomSend")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomSend.verify(message.roomSend);
                    if (error)
                        return "roomSend." + error;
                }
            }
            if (message.inventoryAdd != null && message.hasOwnProperty("inventoryAdd")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryAdd.verify(message.inventoryAdd);
                    if (error)
                        return "inventoryAdd." + error;
                }
            }
            if (message.inventoryConsume != null && message.hasOwnProperty("inventoryConsume")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryConsume.verify(message.inventoryConsume);
                    if (error)
                        return "inventoryConsume." + error;
                }
            }
            if (message.inventoryDrop != null && message.hasOwnProperty("inventoryDrop")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryDrop.verify(message.inventoryDrop);
                    if (error)
                        return "inventoryDrop." + error;
                }
            }
            if (message.inventoryGet != null && message.hasOwnProperty("inventoryGet")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryGet.verify(message.inventoryGet);
                    if (error)
                        return "inventoryGet." + error;
                }
            }
            if (message.msgJson != null && message.hasOwnProperty("msgJson"))
                if (!$util.isString(message.msgJson))
                    return "msgJson: string expected";
            return null;
        };

        /**
         * Creates a WSReq message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.WSReq
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.WSReq} WSReq
         */
        WSReq.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.WSReq)
                return object;
            let message = new $root.ws.WSReq();
            if (object.seq != null)
                if ($util.Long)
                    (message.seq = $util.Long.fromValue(object.seq)).unsigned = false;
                else if (typeof object.seq === "string")
                    message.seq = parseInt(object.seq, 10);
                else if (typeof object.seq === "number")
                    message.seq = object.seq;
                else if (typeof object.seq === "object")
                    message.seq = new $util.LongBits(object.seq.low >>> 0, object.seq.high >>> 0).toNumber();
            if (object.name != null)
                message.name = String(object.name);
            if (object.proxy != null)
                message.proxy = String(object.proxy);
            if (object.authUid != null)
                message.authUid = String(object.authUid);
            if (object.authUsername != null)
                message.authUsername = String(object.authUsername);
            if (object.accountLogin != null) {
                if (typeof object.accountLogin !== "object")
                    throw TypeError(".ws.WSReq.accountLogin: object expected");
                message.accountLogin = $root.ws.AccountLogin.fromObject(object.accountLogin);
            }
            if (object.accountRegister != null) {
                if (typeof object.accountRegister !== "object")
                    throw TypeError(".ws.WSReq.accountRegister: object expected");
                message.accountRegister = $root.ws.AccountRegister.fromObject(object.accountRegister);
            }
            if (object.chatWorldSend != null) {
                if (typeof object.chatWorldSend !== "object")
                    throw TypeError(".ws.WSReq.chatWorldSend: object expected");
                message.chatWorldSend = $root.ws.ChatWorldSend.fromObject(object.chatWorldSend);
            }
            if (object.chatPrivateSend != null) {
                if (typeof object.chatPrivateSend !== "object")
                    throw TypeError(".ws.WSReq.chatPrivateSend: object expected");
                message.chatPrivateSend = $root.ws.ChatPrivateSend.fromObject(object.chatPrivateSend);
            }
            if (object.chatGroupSend != null) {
                if (typeof object.chatGroupSend !== "object")
                    throw TypeError(".ws.WSReq.chatGroupSend: object expected");
                message.chatGroupSend = $root.ws.ChatGroupSend.fromObject(object.chatGroupSend);
            }
            if (object.groupMemberReq != null) {
                if (typeof object.groupMemberReq !== "object")
                    throw TypeError(".ws.WSReq.groupMemberReq: object expected");
                message.groupMemberReq = $root.ws.GroupMember.fromObject(object.groupMemberReq);
            }
            if (object.roomCreate != null) {
                if (typeof object.roomCreate !== "object")
                    throw TypeError(".ws.WSReq.roomCreate: object expected");
                message.roomCreate = $root.ws.RoomCreate.fromObject(object.roomCreate);
            }
            if (object.roomJoin != null) {
                if (typeof object.roomJoin !== "object")
                    throw TypeError(".ws.WSReq.roomJoin: object expected");
                message.roomJoin = $root.ws.RoomJoin.fromObject(object.roomJoin);
            }
            if (object.roomLeave != null) {
                if (typeof object.roomLeave !== "object")
                    throw TypeError(".ws.WSReq.roomLeave: object expected");
                message.roomLeave = $root.ws.RoomLeave.fromObject(object.roomLeave);
            }
            if (object.roomKick != null) {
                if (typeof object.roomKick !== "object")
                    throw TypeError(".ws.WSReq.roomKick: object expected");
                message.roomKick = $root.ws.RoomKick.fromObject(object.roomKick);
            }
            if (object.roomDelete != null) {
                if (typeof object.roomDelete !== "object")
                    throw TypeError(".ws.WSReq.roomDelete: object expected");
                message.roomDelete = $root.ws.RoomDelete.fromObject(object.roomDelete);
            }
            if (object.roomSend != null) {
                if (typeof object.roomSend !== "object")
                    throw TypeError(".ws.WSReq.roomSend: object expected");
                message.roomSend = $root.ws.RoomSend.fromObject(object.roomSend);
            }
            if (object.inventoryAdd != null) {
                if (typeof object.inventoryAdd !== "object")
                    throw TypeError(".ws.WSReq.inventoryAdd: object expected");
                message.inventoryAdd = $root.ws.InventoryAdd.fromObject(object.inventoryAdd);
            }
            if (object.inventoryConsume != null) {
                if (typeof object.inventoryConsume !== "object")
                    throw TypeError(".ws.WSReq.inventoryConsume: object expected");
                message.inventoryConsume = $root.ws.InventoryConsume.fromObject(object.inventoryConsume);
            }
            if (object.inventoryDrop != null) {
                if (typeof object.inventoryDrop !== "object")
                    throw TypeError(".ws.WSReq.inventoryDrop: object expected");
                message.inventoryDrop = $root.ws.InventoryDrop.fromObject(object.inventoryDrop);
            }
            if (object.inventoryGet != null) {
                if (typeof object.inventoryGet !== "object")
                    throw TypeError(".ws.WSReq.inventoryGet: object expected");
                message.inventoryGet = $root.ws.InventoryGet.fromObject(object.inventoryGet);
            }
            if (object.msgJson != null)
                message.msgJson = String(object.msgJson);
            return message;
        };

        /**
         * Creates a plain object from a WSReq message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.WSReq
         * @static
         * @param {ws.WSReq} message WSReq
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        WSReq.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.seq = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.seq = options.longs === String ? "0" : 0;
                object.name = "";
                object.proxy = "";
                object.authUid = "";
                object.authUsername = "";
                object.msgJson = "";
            }
            if (message.seq != null && message.hasOwnProperty("seq"))
                if (typeof message.seq === "number")
                    object.seq = options.longs === String ? String(message.seq) : message.seq;
                else
                    object.seq = options.longs === String ? $util.Long.prototype.toString.call(message.seq) : options.longs === Number ? new $util.LongBits(message.seq.low >>> 0, message.seq.high >>> 0).toNumber() : message.seq;
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            if (message.proxy != null && message.hasOwnProperty("proxy"))
                object.proxy = message.proxy;
            if (message.authUid != null && message.hasOwnProperty("authUid"))
                object.authUid = message.authUid;
            if (message.authUsername != null && message.hasOwnProperty("authUsername"))
                object.authUsername = message.authUsername;
            if (message.accountLogin != null && message.hasOwnProperty("accountLogin")) {
                object.accountLogin = $root.ws.AccountLogin.toObject(message.accountLogin, options);
                if (options.oneofs)
                    object.msg = "accountLogin";
            }
            if (message.accountRegister != null && message.hasOwnProperty("accountRegister")) {
                object.accountRegister = $root.ws.AccountRegister.toObject(message.accountRegister, options);
                if (options.oneofs)
                    object.msg = "accountRegister";
            }
            if (message.chatWorldSend != null && message.hasOwnProperty("chatWorldSend")) {
                object.chatWorldSend = $root.ws.ChatWorldSend.toObject(message.chatWorldSend, options);
                if (options.oneofs)
                    object.msg = "chatWorldSend";
            }
            if (message.chatPrivateSend != null && message.hasOwnProperty("chatPrivateSend")) {
                object.chatPrivateSend = $root.ws.ChatPrivateSend.toObject(message.chatPrivateSend, options);
                if (options.oneofs)
                    object.msg = "chatPrivateSend";
            }
            if (message.chatGroupSend != null && message.hasOwnProperty("chatGroupSend")) {
                object.chatGroupSend = $root.ws.ChatGroupSend.toObject(message.chatGroupSend, options);
                if (options.oneofs)
                    object.msg = "chatGroupSend";
            }
            if (message.groupMemberReq != null && message.hasOwnProperty("groupMemberReq")) {
                object.groupMemberReq = $root.ws.GroupMember.toObject(message.groupMemberReq, options);
                if (options.oneofs)
                    object.msg = "groupMemberReq";
            }
            if (message.roomCreate != null && message.hasOwnProperty("roomCreate")) {
                object.roomCreate = $root.ws.RoomCreate.toObject(message.roomCreate, options);
                if (options.oneofs)
                    object.msg = "roomCreate";
            }
            if (message.roomJoin != null && message.hasOwnProperty("roomJoin")) {
                object.roomJoin = $root.ws.RoomJoin.toObject(message.roomJoin, options);
                if (options.oneofs)
                    object.msg = "roomJoin";
            }
            if (message.roomLeave != null && message.hasOwnProperty("roomLeave")) {
                object.roomLeave = $root.ws.RoomLeave.toObject(message.roomLeave, options);
                if (options.oneofs)
                    object.msg = "roomLeave";
            }
            if (message.roomKick != null && message.hasOwnProperty("roomKick")) {
                object.roomKick = $root.ws.RoomKick.toObject(message.roomKick, options);
                if (options.oneofs)
                    object.msg = "roomKick";
            }
            if (message.roomDelete != null && message.hasOwnProperty("roomDelete")) {
                object.roomDelete = $root.ws.RoomDelete.toObject(message.roomDelete, options);
                if (options.oneofs)
                    object.msg = "roomDelete";
            }
            if (message.roomSend != null && message.hasOwnProperty("roomSend")) {
                object.roomSend = $root.ws.RoomSend.toObject(message.roomSend, options);
                if (options.oneofs)
                    object.msg = "roomSend";
            }
            if (message.inventoryAdd != null && message.hasOwnProperty("inventoryAdd")) {
                object.inventoryAdd = $root.ws.InventoryAdd.toObject(message.inventoryAdd, options);
                if (options.oneofs)
                    object.msg = "inventoryAdd";
            }
            if (message.inventoryConsume != null && message.hasOwnProperty("inventoryConsume")) {
                object.inventoryConsume = $root.ws.InventoryConsume.toObject(message.inventoryConsume, options);
                if (options.oneofs)
                    object.msg = "inventoryConsume";
            }
            if (message.inventoryDrop != null && message.hasOwnProperty("inventoryDrop")) {
                object.inventoryDrop = $root.ws.InventoryDrop.toObject(message.inventoryDrop, options);
                if (options.oneofs)
                    object.msg = "inventoryDrop";
            }
            if (message.inventoryGet != null && message.hasOwnProperty("inventoryGet")) {
                object.inventoryGet = $root.ws.InventoryGet.toObject(message.inventoryGet, options);
                if (options.oneofs)
                    object.msg = "inventoryGet";
            }
            if (message.msgJson != null && message.hasOwnProperty("msgJson"))
                object.msgJson = message.msgJson;
            return object;
        };

        /**
         * Converts this WSReq to JSON.
         * @function toJSON
         * @memberof ws.WSReq
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        WSReq.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return WSReq;
    })();

    ws.WSResp = (function() {

        /**
         * Properties of a WSResp.
         * @memberof ws
         * @interface IWSResp
         * @property {number|Long|null} [seq] WSResp seq
         * @property {string|null} [name] WSResp name
         * @property {number|null} [code] WSResp code
         * @property {ws.IAccountLoginResp|null} [accountLoginResp] WSResp accountLoginResp
         * @property {ws.IAccountRegisterResp|null} [accountRegisterResp] WSResp accountRegisterResp
         * @property {ws.IChatWorldNew|null} [chatWorldNew] WSResp chatWorldNew
         * @property {ws.IChatPrivateNew|null} [chatPrivateNew] WSResp chatPrivateNew
         * @property {ws.IChatGroupNew|null} [chatGroupNew] WSResp chatGroupNew
         * @property {ws.IGroupMemberResp|null} [groupMemberResp] WSResp groupMemberResp
         * @property {ws.IChatOnlineResp|null} [chatOnlineResp] WSResp chatOnlineResp
         * @property {ws.IRoomCreateResp|null} [roomCreateResp] WSResp roomCreateResp
         * @property {ws.IRoomJoinResp|null} [roomJoinResp] WSResp roomJoinResp
         * @property {ws.IRoomChatNew|null} [roomChatNew] WSResp roomChatNew
         * @property {ws.IRoomNotice|null} [roomNotice] WSResp roomNotice
         * @property {ws.ISystemBroadcast|null} [systemBroadcast] WSResp systemBroadcast
         * @property {ws.ISystemPush|null} [systemPush] WSResp systemPush
         * @property {ws.IInventoryAddResp|null} [inventoryAddResp] WSResp inventoryAddResp
         * @property {ws.IInventoryConsumeResp|null} [inventoryConsumeResp] WSResp inventoryConsumeResp
         * @property {ws.IInventoryDropResp|null} [inventoryDropResp] WSResp inventoryDropResp
         * @property {ws.IInventoryGetResp|null} [inventoryGetResp] WSResp inventoryGetResp
         * @property {string|null} [msgJson] WSResp msgJson
         */

        /**
         * Constructs a new WSResp.
         * @memberof ws
         * @classdesc Represents a WSResp.
         * @implements IWSResp
         * @constructor
         * @param {ws.IWSResp=} [properties] Properties to set
         */
        function WSResp(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * WSResp seq.
         * @member {number|Long} seq
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.seq = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * WSResp name.
         * @member {string} name
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.name = "";

        /**
         * WSResp code.
         * @member {number} code
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.code = 0;

        /**
         * WSResp accountLoginResp.
         * @member {ws.IAccountLoginResp|null|undefined} accountLoginResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.accountLoginResp = null;

        /**
         * WSResp accountRegisterResp.
         * @member {ws.IAccountRegisterResp|null|undefined} accountRegisterResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.accountRegisterResp = null;

        /**
         * WSResp chatWorldNew.
         * @member {ws.IChatWorldNew|null|undefined} chatWorldNew
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.chatWorldNew = null;

        /**
         * WSResp chatPrivateNew.
         * @member {ws.IChatPrivateNew|null|undefined} chatPrivateNew
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.chatPrivateNew = null;

        /**
         * WSResp chatGroupNew.
         * @member {ws.IChatGroupNew|null|undefined} chatGroupNew
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.chatGroupNew = null;

        /**
         * WSResp groupMemberResp.
         * @member {ws.IGroupMemberResp|null|undefined} groupMemberResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.groupMemberResp = null;

        /**
         * WSResp chatOnlineResp.
         * @member {ws.IChatOnlineResp|null|undefined} chatOnlineResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.chatOnlineResp = null;

        /**
         * WSResp roomCreateResp.
         * @member {ws.IRoomCreateResp|null|undefined} roomCreateResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.roomCreateResp = null;

        /**
         * WSResp roomJoinResp.
         * @member {ws.IRoomJoinResp|null|undefined} roomJoinResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.roomJoinResp = null;

        /**
         * WSResp roomChatNew.
         * @member {ws.IRoomChatNew|null|undefined} roomChatNew
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.roomChatNew = null;

        /**
         * WSResp roomNotice.
         * @member {ws.IRoomNotice|null|undefined} roomNotice
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.roomNotice = null;

        /**
         * WSResp systemBroadcast.
         * @member {ws.ISystemBroadcast|null|undefined} systemBroadcast
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.systemBroadcast = null;

        /**
         * WSResp systemPush.
         * @member {ws.ISystemPush|null|undefined} systemPush
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.systemPush = null;

        /**
         * WSResp inventoryAddResp.
         * @member {ws.IInventoryAddResp|null|undefined} inventoryAddResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.inventoryAddResp = null;

        /**
         * WSResp inventoryConsumeResp.
         * @member {ws.IInventoryConsumeResp|null|undefined} inventoryConsumeResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.inventoryConsumeResp = null;

        /**
         * WSResp inventoryDropResp.
         * @member {ws.IInventoryDropResp|null|undefined} inventoryDropResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.inventoryDropResp = null;

        /**
         * WSResp inventoryGetResp.
         * @member {ws.IInventoryGetResp|null|undefined} inventoryGetResp
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.inventoryGetResp = null;

        /**
         * WSResp msgJson.
         * @member {string} msgJson
         * @memberof ws.WSResp
         * @instance
         */
        WSResp.prototype.msgJson = "";

        // OneOf field names bound to virtual getters and setters
        let $oneOfFields;

        /**
         * WSResp msg.
         * @member {"accountLoginResp"|"accountRegisterResp"|"chatWorldNew"|"chatPrivateNew"|"chatGroupNew"|"groupMemberResp"|"chatOnlineResp"|"roomCreateResp"|"roomJoinResp"|"roomChatNew"|"roomNotice"|"systemBroadcast"|"systemPush"|"inventoryAddResp"|"inventoryConsumeResp"|"inventoryDropResp"|"inventoryGetResp"|undefined} msg
         * @memberof ws.WSResp
         * @instance
         */
        Object.defineProperty(WSResp.prototype, "msg", {
            get: $util.oneOfGetter($oneOfFields = ["accountLoginResp", "accountRegisterResp", "chatWorldNew", "chatPrivateNew", "chatGroupNew", "groupMemberResp", "chatOnlineResp", "roomCreateResp", "roomJoinResp", "roomChatNew", "roomNotice", "systemBroadcast", "systemPush", "inventoryAddResp", "inventoryConsumeResp", "inventoryDropResp", "inventoryGetResp"]),
            set: $util.oneOfSetter($oneOfFields)
        });

        /**
         * Creates a new WSResp instance using the specified properties.
         * @function create
         * @memberof ws.WSResp
         * @static
         * @param {ws.IWSResp=} [properties] Properties to set
         * @returns {ws.WSResp} WSResp instance
         */
        WSResp.create = function create(properties) {
            return new WSResp(properties);
        };

        /**
         * Encodes the specified WSResp message. Does not implicitly {@link ws.WSResp.verify|verify} messages.
         * @function encode
         * @memberof ws.WSResp
         * @static
         * @param {ws.IWSResp} message WSResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        WSResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.seq != null && Object.hasOwnProperty.call(message, "seq"))
                writer.uint32(/* id 1, wireType 0 =*/8).int64(message.seq);
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
            if (message.code != null && Object.hasOwnProperty.call(message, "code"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.code);
            if (message.accountLoginResp != null && Object.hasOwnProperty.call(message, "accountLoginResp"))
                $root.ws.AccountLoginResp.encode(message.accountLoginResp, writer.uint32(/* id 10, wireType 2 =*/82).fork()).ldelim();
            if (message.accountRegisterResp != null && Object.hasOwnProperty.call(message, "accountRegisterResp"))
                $root.ws.AccountRegisterResp.encode(message.accountRegisterResp, writer.uint32(/* id 11, wireType 2 =*/90).fork()).ldelim();
            if (message.chatWorldNew != null && Object.hasOwnProperty.call(message, "chatWorldNew"))
                $root.ws.ChatWorldNew.encode(message.chatWorldNew, writer.uint32(/* id 20, wireType 2 =*/162).fork()).ldelim();
            if (message.chatPrivateNew != null && Object.hasOwnProperty.call(message, "chatPrivateNew"))
                $root.ws.ChatPrivateNew.encode(message.chatPrivateNew, writer.uint32(/* id 21, wireType 2 =*/170).fork()).ldelim();
            if (message.chatGroupNew != null && Object.hasOwnProperty.call(message, "chatGroupNew"))
                $root.ws.ChatGroupNew.encode(message.chatGroupNew, writer.uint32(/* id 22, wireType 2 =*/178).fork()).ldelim();
            if (message.groupMemberResp != null && Object.hasOwnProperty.call(message, "groupMemberResp"))
                $root.ws.GroupMemberResp.encode(message.groupMemberResp, writer.uint32(/* id 23, wireType 2 =*/186).fork()).ldelim();
            if (message.chatOnlineResp != null && Object.hasOwnProperty.call(message, "chatOnlineResp"))
                $root.ws.ChatOnlineResp.encode(message.chatOnlineResp, writer.uint32(/* id 24, wireType 2 =*/194).fork()).ldelim();
            if (message.roomCreateResp != null && Object.hasOwnProperty.call(message, "roomCreateResp"))
                $root.ws.RoomCreateResp.encode(message.roomCreateResp, writer.uint32(/* id 30, wireType 2 =*/242).fork()).ldelim();
            if (message.roomJoinResp != null && Object.hasOwnProperty.call(message, "roomJoinResp"))
                $root.ws.RoomJoinResp.encode(message.roomJoinResp, writer.uint32(/* id 31, wireType 2 =*/250).fork()).ldelim();
            if (message.roomChatNew != null && Object.hasOwnProperty.call(message, "roomChatNew"))
                $root.ws.RoomChatNew.encode(message.roomChatNew, writer.uint32(/* id 32, wireType 2 =*/258).fork()).ldelim();
            if (message.roomNotice != null && Object.hasOwnProperty.call(message, "roomNotice"))
                $root.ws.RoomNotice.encode(message.roomNotice, writer.uint32(/* id 33, wireType 2 =*/266).fork()).ldelim();
            if (message.systemBroadcast != null && Object.hasOwnProperty.call(message, "systemBroadcast"))
                $root.ws.SystemBroadcast.encode(message.systemBroadcast, writer.uint32(/* id 40, wireType 2 =*/322).fork()).ldelim();
            if (message.systemPush != null && Object.hasOwnProperty.call(message, "systemPush"))
                $root.ws.SystemPush.encode(message.systemPush, writer.uint32(/* id 41, wireType 2 =*/330).fork()).ldelim();
            if (message.inventoryAddResp != null && Object.hasOwnProperty.call(message, "inventoryAddResp"))
                $root.ws.InventoryAddResp.encode(message.inventoryAddResp, writer.uint32(/* id 50, wireType 2 =*/402).fork()).ldelim();
            if (message.inventoryConsumeResp != null && Object.hasOwnProperty.call(message, "inventoryConsumeResp"))
                $root.ws.InventoryConsumeResp.encode(message.inventoryConsumeResp, writer.uint32(/* id 51, wireType 2 =*/410).fork()).ldelim();
            if (message.inventoryDropResp != null && Object.hasOwnProperty.call(message, "inventoryDropResp"))
                $root.ws.InventoryDropResp.encode(message.inventoryDropResp, writer.uint32(/* id 52, wireType 2 =*/418).fork()).ldelim();
            if (message.inventoryGetResp != null && Object.hasOwnProperty.call(message, "inventoryGetResp"))
                $root.ws.InventoryGetResp.encode(message.inventoryGetResp, writer.uint32(/* id 53, wireType 2 =*/426).fork()).ldelim();
            if (message.msgJson != null && Object.hasOwnProperty.call(message, "msgJson"))
                writer.uint32(/* id 100, wireType 2 =*/802).string(message.msgJson);
            return writer;
        };

        /**
         * Encodes the specified WSResp message, length delimited. Does not implicitly {@link ws.WSResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.WSResp
         * @static
         * @param {ws.IWSResp} message WSResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        WSResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a WSResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.WSResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.WSResp} WSResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        WSResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.WSResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.seq = reader.int64();
                    break;
                case 2:
                    message.name = reader.string();
                    break;
                case 3:
                    message.code = reader.int32();
                    break;
                case 10:
                    message.accountLoginResp = $root.ws.AccountLoginResp.decode(reader, reader.uint32());
                    break;
                case 11:
                    message.accountRegisterResp = $root.ws.AccountRegisterResp.decode(reader, reader.uint32());
                    break;
                case 20:
                    message.chatWorldNew = $root.ws.ChatWorldNew.decode(reader, reader.uint32());
                    break;
                case 21:
                    message.chatPrivateNew = $root.ws.ChatPrivateNew.decode(reader, reader.uint32());
                    break;
                case 22:
                    message.chatGroupNew = $root.ws.ChatGroupNew.decode(reader, reader.uint32());
                    break;
                case 23:
                    message.groupMemberResp = $root.ws.GroupMemberResp.decode(reader, reader.uint32());
                    break;
                case 24:
                    message.chatOnlineResp = $root.ws.ChatOnlineResp.decode(reader, reader.uint32());
                    break;
                case 30:
                    message.roomCreateResp = $root.ws.RoomCreateResp.decode(reader, reader.uint32());
                    break;
                case 31:
                    message.roomJoinResp = $root.ws.RoomJoinResp.decode(reader, reader.uint32());
                    break;
                case 32:
                    message.roomChatNew = $root.ws.RoomChatNew.decode(reader, reader.uint32());
                    break;
                case 33:
                    message.roomNotice = $root.ws.RoomNotice.decode(reader, reader.uint32());
                    break;
                case 40:
                    message.systemBroadcast = $root.ws.SystemBroadcast.decode(reader, reader.uint32());
                    break;
                case 41:
                    message.systemPush = $root.ws.SystemPush.decode(reader, reader.uint32());
                    break;
                case 50:
                    message.inventoryAddResp = $root.ws.InventoryAddResp.decode(reader, reader.uint32());
                    break;
                case 51:
                    message.inventoryConsumeResp = $root.ws.InventoryConsumeResp.decode(reader, reader.uint32());
                    break;
                case 52:
                    message.inventoryDropResp = $root.ws.InventoryDropResp.decode(reader, reader.uint32());
                    break;
                case 53:
                    message.inventoryGetResp = $root.ws.InventoryGetResp.decode(reader, reader.uint32());
                    break;
                case 100:
                    message.msgJson = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a WSResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.WSResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.WSResp} WSResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        WSResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a WSResp message.
         * @function verify
         * @memberof ws.WSResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        WSResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            let properties = {};
            if (message.seq != null && message.hasOwnProperty("seq"))
                if (!$util.isInteger(message.seq) && !(message.seq && $util.isInteger(message.seq.low) && $util.isInteger(message.seq.high)))
                    return "seq: integer|Long expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            if (message.code != null && message.hasOwnProperty("code"))
                if (!$util.isInteger(message.code))
                    return "code: integer expected";
            if (message.accountLoginResp != null && message.hasOwnProperty("accountLoginResp")) {
                properties.msg = 1;
                {
                    let error = $root.ws.AccountLoginResp.verify(message.accountLoginResp);
                    if (error)
                        return "accountLoginResp." + error;
                }
            }
            if (message.accountRegisterResp != null && message.hasOwnProperty("accountRegisterResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.AccountRegisterResp.verify(message.accountRegisterResp);
                    if (error)
                        return "accountRegisterResp." + error;
                }
            }
            if (message.chatWorldNew != null && message.hasOwnProperty("chatWorldNew")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatWorldNew.verify(message.chatWorldNew);
                    if (error)
                        return "chatWorldNew." + error;
                }
            }
            if (message.chatPrivateNew != null && message.hasOwnProperty("chatPrivateNew")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatPrivateNew.verify(message.chatPrivateNew);
                    if (error)
                        return "chatPrivateNew." + error;
                }
            }
            if (message.chatGroupNew != null && message.hasOwnProperty("chatGroupNew")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatGroupNew.verify(message.chatGroupNew);
                    if (error)
                        return "chatGroupNew." + error;
                }
            }
            if (message.groupMemberResp != null && message.hasOwnProperty("groupMemberResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.GroupMemberResp.verify(message.groupMemberResp);
                    if (error)
                        return "groupMemberResp." + error;
                }
            }
            if (message.chatOnlineResp != null && message.hasOwnProperty("chatOnlineResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.ChatOnlineResp.verify(message.chatOnlineResp);
                    if (error)
                        return "chatOnlineResp." + error;
                }
            }
            if (message.roomCreateResp != null && message.hasOwnProperty("roomCreateResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomCreateResp.verify(message.roomCreateResp);
                    if (error)
                        return "roomCreateResp." + error;
                }
            }
            if (message.roomJoinResp != null && message.hasOwnProperty("roomJoinResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomJoinResp.verify(message.roomJoinResp);
                    if (error)
                        return "roomJoinResp." + error;
                }
            }
            if (message.roomChatNew != null && message.hasOwnProperty("roomChatNew")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomChatNew.verify(message.roomChatNew);
                    if (error)
                        return "roomChatNew." + error;
                }
            }
            if (message.roomNotice != null && message.hasOwnProperty("roomNotice")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.RoomNotice.verify(message.roomNotice);
                    if (error)
                        return "roomNotice." + error;
                }
            }
            if (message.systemBroadcast != null && message.hasOwnProperty("systemBroadcast")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.SystemBroadcast.verify(message.systemBroadcast);
                    if (error)
                        return "systemBroadcast." + error;
                }
            }
            if (message.systemPush != null && message.hasOwnProperty("systemPush")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.SystemPush.verify(message.systemPush);
                    if (error)
                        return "systemPush." + error;
                }
            }
            if (message.inventoryAddResp != null && message.hasOwnProperty("inventoryAddResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryAddResp.verify(message.inventoryAddResp);
                    if (error)
                        return "inventoryAddResp." + error;
                }
            }
            if (message.inventoryConsumeResp != null && message.hasOwnProperty("inventoryConsumeResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryConsumeResp.verify(message.inventoryConsumeResp);
                    if (error)
                        return "inventoryConsumeResp." + error;
                }
            }
            if (message.inventoryDropResp != null && message.hasOwnProperty("inventoryDropResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryDropResp.verify(message.inventoryDropResp);
                    if (error)
                        return "inventoryDropResp." + error;
                }
            }
            if (message.inventoryGetResp != null && message.hasOwnProperty("inventoryGetResp")) {
                if (properties.msg === 1)
                    return "msg: multiple values";
                properties.msg = 1;
                {
                    let error = $root.ws.InventoryGetResp.verify(message.inventoryGetResp);
                    if (error)
                        return "inventoryGetResp." + error;
                }
            }
            if (message.msgJson != null && message.hasOwnProperty("msgJson"))
                if (!$util.isString(message.msgJson))
                    return "msgJson: string expected";
            return null;
        };

        /**
         * Creates a WSResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.WSResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.WSResp} WSResp
         */
        WSResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.WSResp)
                return object;
            let message = new $root.ws.WSResp();
            if (object.seq != null)
                if ($util.Long)
                    (message.seq = $util.Long.fromValue(object.seq)).unsigned = false;
                else if (typeof object.seq === "string")
                    message.seq = parseInt(object.seq, 10);
                else if (typeof object.seq === "number")
                    message.seq = object.seq;
                else if (typeof object.seq === "object")
                    message.seq = new $util.LongBits(object.seq.low >>> 0, object.seq.high >>> 0).toNumber();
            if (object.name != null)
                message.name = String(object.name);
            if (object.code != null)
                message.code = object.code | 0;
            if (object.accountLoginResp != null) {
                if (typeof object.accountLoginResp !== "object")
                    throw TypeError(".ws.WSResp.accountLoginResp: object expected");
                message.accountLoginResp = $root.ws.AccountLoginResp.fromObject(object.accountLoginResp);
            }
            if (object.accountRegisterResp != null) {
                if (typeof object.accountRegisterResp !== "object")
                    throw TypeError(".ws.WSResp.accountRegisterResp: object expected");
                message.accountRegisterResp = $root.ws.AccountRegisterResp.fromObject(object.accountRegisterResp);
            }
            if (object.chatWorldNew != null) {
                if (typeof object.chatWorldNew !== "object")
                    throw TypeError(".ws.WSResp.chatWorldNew: object expected");
                message.chatWorldNew = $root.ws.ChatWorldNew.fromObject(object.chatWorldNew);
            }
            if (object.chatPrivateNew != null) {
                if (typeof object.chatPrivateNew !== "object")
                    throw TypeError(".ws.WSResp.chatPrivateNew: object expected");
                message.chatPrivateNew = $root.ws.ChatPrivateNew.fromObject(object.chatPrivateNew);
            }
            if (object.chatGroupNew != null) {
                if (typeof object.chatGroupNew !== "object")
                    throw TypeError(".ws.WSResp.chatGroupNew: object expected");
                message.chatGroupNew = $root.ws.ChatGroupNew.fromObject(object.chatGroupNew);
            }
            if (object.groupMemberResp != null) {
                if (typeof object.groupMemberResp !== "object")
                    throw TypeError(".ws.WSResp.groupMemberResp: object expected");
                message.groupMemberResp = $root.ws.GroupMemberResp.fromObject(object.groupMemberResp);
            }
            if (object.chatOnlineResp != null) {
                if (typeof object.chatOnlineResp !== "object")
                    throw TypeError(".ws.WSResp.chatOnlineResp: object expected");
                message.chatOnlineResp = $root.ws.ChatOnlineResp.fromObject(object.chatOnlineResp);
            }
            if (object.roomCreateResp != null) {
                if (typeof object.roomCreateResp !== "object")
                    throw TypeError(".ws.WSResp.roomCreateResp: object expected");
                message.roomCreateResp = $root.ws.RoomCreateResp.fromObject(object.roomCreateResp);
            }
            if (object.roomJoinResp != null) {
                if (typeof object.roomJoinResp !== "object")
                    throw TypeError(".ws.WSResp.roomJoinResp: object expected");
                message.roomJoinResp = $root.ws.RoomJoinResp.fromObject(object.roomJoinResp);
            }
            if (object.roomChatNew != null) {
                if (typeof object.roomChatNew !== "object")
                    throw TypeError(".ws.WSResp.roomChatNew: object expected");
                message.roomChatNew = $root.ws.RoomChatNew.fromObject(object.roomChatNew);
            }
            if (object.roomNotice != null) {
                if (typeof object.roomNotice !== "object")
                    throw TypeError(".ws.WSResp.roomNotice: object expected");
                message.roomNotice = $root.ws.RoomNotice.fromObject(object.roomNotice);
            }
            if (object.systemBroadcast != null) {
                if (typeof object.systemBroadcast !== "object")
                    throw TypeError(".ws.WSResp.systemBroadcast: object expected");
                message.systemBroadcast = $root.ws.SystemBroadcast.fromObject(object.systemBroadcast);
            }
            if (object.systemPush != null) {
                if (typeof object.systemPush !== "object")
                    throw TypeError(".ws.WSResp.systemPush: object expected");
                message.systemPush = $root.ws.SystemPush.fromObject(object.systemPush);
            }
            if (object.inventoryAddResp != null) {
                if (typeof object.inventoryAddResp !== "object")
                    throw TypeError(".ws.WSResp.inventoryAddResp: object expected");
                message.inventoryAddResp = $root.ws.InventoryAddResp.fromObject(object.inventoryAddResp);
            }
            if (object.inventoryConsumeResp != null) {
                if (typeof object.inventoryConsumeResp !== "object")
                    throw TypeError(".ws.WSResp.inventoryConsumeResp: object expected");
                message.inventoryConsumeResp = $root.ws.InventoryConsumeResp.fromObject(object.inventoryConsumeResp);
            }
            if (object.inventoryDropResp != null) {
                if (typeof object.inventoryDropResp !== "object")
                    throw TypeError(".ws.WSResp.inventoryDropResp: object expected");
                message.inventoryDropResp = $root.ws.InventoryDropResp.fromObject(object.inventoryDropResp);
            }
            if (object.inventoryGetResp != null) {
                if (typeof object.inventoryGetResp !== "object")
                    throw TypeError(".ws.WSResp.inventoryGetResp: object expected");
                message.inventoryGetResp = $root.ws.InventoryGetResp.fromObject(object.inventoryGetResp);
            }
            if (object.msgJson != null)
                message.msgJson = String(object.msgJson);
            return message;
        };

        /**
         * Creates a plain object from a WSResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.WSResp
         * @static
         * @param {ws.WSResp} message WSResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        WSResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.seq = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.seq = options.longs === String ? "0" : 0;
                object.name = "";
                object.code = 0;
                object.msgJson = "";
            }
            if (message.seq != null && message.hasOwnProperty("seq"))
                if (typeof message.seq === "number")
                    object.seq = options.longs === String ? String(message.seq) : message.seq;
                else
                    object.seq = options.longs === String ? $util.Long.prototype.toString.call(message.seq) : options.longs === Number ? new $util.LongBits(message.seq.low >>> 0, message.seq.high >>> 0).toNumber() : message.seq;
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            if (message.code != null && message.hasOwnProperty("code"))
                object.code = message.code;
            if (message.accountLoginResp != null && message.hasOwnProperty("accountLoginResp")) {
                object.accountLoginResp = $root.ws.AccountLoginResp.toObject(message.accountLoginResp, options);
                if (options.oneofs)
                    object.msg = "accountLoginResp";
            }
            if (message.accountRegisterResp != null && message.hasOwnProperty("accountRegisterResp")) {
                object.accountRegisterResp = $root.ws.AccountRegisterResp.toObject(message.accountRegisterResp, options);
                if (options.oneofs)
                    object.msg = "accountRegisterResp";
            }
            if (message.chatWorldNew != null && message.hasOwnProperty("chatWorldNew")) {
                object.chatWorldNew = $root.ws.ChatWorldNew.toObject(message.chatWorldNew, options);
                if (options.oneofs)
                    object.msg = "chatWorldNew";
            }
            if (message.chatPrivateNew != null && message.hasOwnProperty("chatPrivateNew")) {
                object.chatPrivateNew = $root.ws.ChatPrivateNew.toObject(message.chatPrivateNew, options);
                if (options.oneofs)
                    object.msg = "chatPrivateNew";
            }
            if (message.chatGroupNew != null && message.hasOwnProperty("chatGroupNew")) {
                object.chatGroupNew = $root.ws.ChatGroupNew.toObject(message.chatGroupNew, options);
                if (options.oneofs)
                    object.msg = "chatGroupNew";
            }
            if (message.groupMemberResp != null && message.hasOwnProperty("groupMemberResp")) {
                object.groupMemberResp = $root.ws.GroupMemberResp.toObject(message.groupMemberResp, options);
                if (options.oneofs)
                    object.msg = "groupMemberResp";
            }
            if (message.chatOnlineResp != null && message.hasOwnProperty("chatOnlineResp")) {
                object.chatOnlineResp = $root.ws.ChatOnlineResp.toObject(message.chatOnlineResp, options);
                if (options.oneofs)
                    object.msg = "chatOnlineResp";
            }
            if (message.roomCreateResp != null && message.hasOwnProperty("roomCreateResp")) {
                object.roomCreateResp = $root.ws.RoomCreateResp.toObject(message.roomCreateResp, options);
                if (options.oneofs)
                    object.msg = "roomCreateResp";
            }
            if (message.roomJoinResp != null && message.hasOwnProperty("roomJoinResp")) {
                object.roomJoinResp = $root.ws.RoomJoinResp.toObject(message.roomJoinResp, options);
                if (options.oneofs)
                    object.msg = "roomJoinResp";
            }
            if (message.roomChatNew != null && message.hasOwnProperty("roomChatNew")) {
                object.roomChatNew = $root.ws.RoomChatNew.toObject(message.roomChatNew, options);
                if (options.oneofs)
                    object.msg = "roomChatNew";
            }
            if (message.roomNotice != null && message.hasOwnProperty("roomNotice")) {
                object.roomNotice = $root.ws.RoomNotice.toObject(message.roomNotice, options);
                if (options.oneofs)
                    object.msg = "roomNotice";
            }
            if (message.systemBroadcast != null && message.hasOwnProperty("systemBroadcast")) {
                object.systemBroadcast = $root.ws.SystemBroadcast.toObject(message.systemBroadcast, options);
                if (options.oneofs)
                    object.msg = "systemBroadcast";
            }
            if (message.systemPush != null && message.hasOwnProperty("systemPush")) {
                object.systemPush = $root.ws.SystemPush.toObject(message.systemPush, options);
                if (options.oneofs)
                    object.msg = "systemPush";
            }
            if (message.inventoryAddResp != null && message.hasOwnProperty("inventoryAddResp")) {
                object.inventoryAddResp = $root.ws.InventoryAddResp.toObject(message.inventoryAddResp, options);
                if (options.oneofs)
                    object.msg = "inventoryAddResp";
            }
            if (message.inventoryConsumeResp != null && message.hasOwnProperty("inventoryConsumeResp")) {
                object.inventoryConsumeResp = $root.ws.InventoryConsumeResp.toObject(message.inventoryConsumeResp, options);
                if (options.oneofs)
                    object.msg = "inventoryConsumeResp";
            }
            if (message.inventoryDropResp != null && message.hasOwnProperty("inventoryDropResp")) {
                object.inventoryDropResp = $root.ws.InventoryDropResp.toObject(message.inventoryDropResp, options);
                if (options.oneofs)
                    object.msg = "inventoryDropResp";
            }
            if (message.inventoryGetResp != null && message.hasOwnProperty("inventoryGetResp")) {
                object.inventoryGetResp = $root.ws.InventoryGetResp.toObject(message.inventoryGetResp, options);
                if (options.oneofs)
                    object.msg = "inventoryGetResp";
            }
            if (message.msgJson != null && message.hasOwnProperty("msgJson"))
                object.msgJson = message.msgJson;
            return object;
        };

        /**
         * Converts this WSResp to JSON.
         * @function toJSON
         * @memberof ws.WSResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        WSResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return WSResp;
    })();

    ws.AccountLogin = (function() {

        /**
         * Properties of an AccountLogin.
         * @memberof ws
         * @interface IAccountLogin
         * @property {string|null} [username] AccountLogin username
         * @property {string|null} [password] AccountLogin password
         */

        /**
         * Constructs a new AccountLogin.
         * @memberof ws
         * @classdesc Represents an AccountLogin.
         * @implements IAccountLogin
         * @constructor
         * @param {ws.IAccountLogin=} [properties] Properties to set
         */
        function AccountLogin(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * AccountLogin username.
         * @member {string} username
         * @memberof ws.AccountLogin
         * @instance
         */
        AccountLogin.prototype.username = "";

        /**
         * AccountLogin password.
         * @member {string} password
         * @memberof ws.AccountLogin
         * @instance
         */
        AccountLogin.prototype.password = "";

        /**
         * Creates a new AccountLogin instance using the specified properties.
         * @function create
         * @memberof ws.AccountLogin
         * @static
         * @param {ws.IAccountLogin=} [properties] Properties to set
         * @returns {ws.AccountLogin} AccountLogin instance
         */
        AccountLogin.create = function create(properties) {
            return new AccountLogin(properties);
        };

        /**
         * Encodes the specified AccountLogin message. Does not implicitly {@link ws.AccountLogin.verify|verify} messages.
         * @function encode
         * @memberof ws.AccountLogin
         * @static
         * @param {ws.IAccountLogin} message AccountLogin message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountLogin.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.username != null && Object.hasOwnProperty.call(message, "username"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.username);
            if (message.password != null && Object.hasOwnProperty.call(message, "password"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.password);
            return writer;
        };

        /**
         * Encodes the specified AccountLogin message, length delimited. Does not implicitly {@link ws.AccountLogin.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.AccountLogin
         * @static
         * @param {ws.IAccountLogin} message AccountLogin message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountLogin.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an AccountLogin message from the specified reader or buffer.
         * @function decode
         * @memberof ws.AccountLogin
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.AccountLogin} AccountLogin
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountLogin.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.AccountLogin();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.username = reader.string();
                    break;
                case 2:
                    message.password = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an AccountLogin message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.AccountLogin
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.AccountLogin} AccountLogin
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountLogin.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an AccountLogin message.
         * @function verify
         * @memberof ws.AccountLogin
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        AccountLogin.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.username != null && message.hasOwnProperty("username"))
                if (!$util.isString(message.username))
                    return "username: string expected";
            if (message.password != null && message.hasOwnProperty("password"))
                if (!$util.isString(message.password))
                    return "password: string expected";
            return null;
        };

        /**
         * Creates an AccountLogin message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.AccountLogin
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.AccountLogin} AccountLogin
         */
        AccountLogin.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.AccountLogin)
                return object;
            let message = new $root.ws.AccountLogin();
            if (object.username != null)
                message.username = String(object.username);
            if (object.password != null)
                message.password = String(object.password);
            return message;
        };

        /**
         * Creates a plain object from an AccountLogin message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.AccountLogin
         * @static
         * @param {ws.AccountLogin} message AccountLogin
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        AccountLogin.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.username = "";
                object.password = "";
            }
            if (message.username != null && message.hasOwnProperty("username"))
                object.username = message.username;
            if (message.password != null && message.hasOwnProperty("password"))
                object.password = message.password;
            return object;
        };

        /**
         * Converts this AccountLogin to JSON.
         * @function toJSON
         * @memberof ws.AccountLogin
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        AccountLogin.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return AccountLogin;
    })();

    ws.AccountLoginResp = (function() {

        /**
         * Properties of an AccountLoginResp.
         * @memberof ws
         * @interface IAccountLoginResp
         * @property {string|null} [username] AccountLoginResp username
         * @property {string|null} [token] AccountLoginResp token
         * @property {string|null} [uid] AccountLoginResp uid
         * @property {Object.<string,number|Long>|null} [items] AccountLoginResp items
         */

        /**
         * Constructs a new AccountLoginResp.
         * @memberof ws
         * @classdesc Represents an AccountLoginResp.
         * @implements IAccountLoginResp
         * @constructor
         * @param {ws.IAccountLoginResp=} [properties] Properties to set
         */
        function AccountLoginResp(properties) {
            this.items = {};
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * AccountLoginResp username.
         * @member {string} username
         * @memberof ws.AccountLoginResp
         * @instance
         */
        AccountLoginResp.prototype.username = "";

        /**
         * AccountLoginResp token.
         * @member {string} token
         * @memberof ws.AccountLoginResp
         * @instance
         */
        AccountLoginResp.prototype.token = "";

        /**
         * AccountLoginResp uid.
         * @member {string} uid
         * @memberof ws.AccountLoginResp
         * @instance
         */
        AccountLoginResp.prototype.uid = "";

        /**
         * AccountLoginResp items.
         * @member {Object.<string,number|Long>} items
         * @memberof ws.AccountLoginResp
         * @instance
         */
        AccountLoginResp.prototype.items = $util.emptyObject;

        /**
         * Creates a new AccountLoginResp instance using the specified properties.
         * @function create
         * @memberof ws.AccountLoginResp
         * @static
         * @param {ws.IAccountLoginResp=} [properties] Properties to set
         * @returns {ws.AccountLoginResp} AccountLoginResp instance
         */
        AccountLoginResp.create = function create(properties) {
            return new AccountLoginResp(properties);
        };

        /**
         * Encodes the specified AccountLoginResp message. Does not implicitly {@link ws.AccountLoginResp.verify|verify} messages.
         * @function encode
         * @memberof ws.AccountLoginResp
         * @static
         * @param {ws.IAccountLoginResp} message AccountLoginResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountLoginResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.username != null && Object.hasOwnProperty.call(message, "username"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.username);
            if (message.token != null && Object.hasOwnProperty.call(message, "token"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.token);
            if (message.uid != null && Object.hasOwnProperty.call(message, "uid"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.uid);
            if (message.items != null && Object.hasOwnProperty.call(message, "items"))
                for (let keys = Object.keys(message.items), i = 0; i < keys.length; ++i)
                    writer.uint32(/* id 4, wireType 2 =*/34).fork().uint32(/* id 1, wireType 2 =*/10).string(keys[i]).uint32(/* id 2, wireType 0 =*/16).int64(message.items[keys[i]]).ldelim();
            return writer;
        };

        /**
         * Encodes the specified AccountLoginResp message, length delimited. Does not implicitly {@link ws.AccountLoginResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.AccountLoginResp
         * @static
         * @param {ws.IAccountLoginResp} message AccountLoginResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountLoginResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an AccountLoginResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.AccountLoginResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.AccountLoginResp} AccountLoginResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountLoginResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.AccountLoginResp(), key, value;
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.username = reader.string();
                    break;
                case 2:
                    message.token = reader.string();
                    break;
                case 3:
                    message.uid = reader.string();
                    break;
                case 4:
                    if (message.items === $util.emptyObject)
                        message.items = {};
                    let end2 = reader.uint32() + reader.pos;
                    key = "";
                    value = 0;
                    while (reader.pos < end2) {
                        let tag2 = reader.uint32();
                        switch (tag2 >>> 3) {
                        case 1:
                            key = reader.string();
                            break;
                        case 2:
                            value = reader.int64();
                            break;
                        default:
                            reader.skipType(tag2 & 7);
                            break;
                        }
                    }
                    message.items[key] = value;
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an AccountLoginResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.AccountLoginResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.AccountLoginResp} AccountLoginResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountLoginResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an AccountLoginResp message.
         * @function verify
         * @memberof ws.AccountLoginResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        AccountLoginResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.username != null && message.hasOwnProperty("username"))
                if (!$util.isString(message.username))
                    return "username: string expected";
            if (message.token != null && message.hasOwnProperty("token"))
                if (!$util.isString(message.token))
                    return "token: string expected";
            if (message.uid != null && message.hasOwnProperty("uid"))
                if (!$util.isString(message.uid))
                    return "uid: string expected";
            if (message.items != null && message.hasOwnProperty("items")) {
                if (!$util.isObject(message.items))
                    return "items: object expected";
                let key = Object.keys(message.items);
                for (let i = 0; i < key.length; ++i)
                    if (!$util.isInteger(message.items[key[i]]) && !(message.items[key[i]] && $util.isInteger(message.items[key[i]].low) && $util.isInteger(message.items[key[i]].high)))
                        return "items: integer|Long{k:string} expected";
            }
            return null;
        };

        /**
         * Creates an AccountLoginResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.AccountLoginResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.AccountLoginResp} AccountLoginResp
         */
        AccountLoginResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.AccountLoginResp)
                return object;
            let message = new $root.ws.AccountLoginResp();
            if (object.username != null)
                message.username = String(object.username);
            if (object.token != null)
                message.token = String(object.token);
            if (object.uid != null)
                message.uid = String(object.uid);
            if (object.items) {
                if (typeof object.items !== "object")
                    throw TypeError(".ws.AccountLoginResp.items: object expected");
                message.items = {};
                for (let keys = Object.keys(object.items), i = 0; i < keys.length; ++i)
                    if ($util.Long)
                        (message.items[keys[i]] = $util.Long.fromValue(object.items[keys[i]])).unsigned = false;
                    else if (typeof object.items[keys[i]] === "string")
                        message.items[keys[i]] = parseInt(object.items[keys[i]], 10);
                    else if (typeof object.items[keys[i]] === "number")
                        message.items[keys[i]] = object.items[keys[i]];
                    else if (typeof object.items[keys[i]] === "object")
                        message.items[keys[i]] = new $util.LongBits(object.items[keys[i]].low >>> 0, object.items[keys[i]].high >>> 0).toNumber();
            }
            return message;
        };

        /**
         * Creates a plain object from an AccountLoginResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.AccountLoginResp
         * @static
         * @param {ws.AccountLoginResp} message AccountLoginResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        AccountLoginResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.objects || options.defaults)
                object.items = {};
            if (options.defaults) {
                object.username = "";
                object.token = "";
                object.uid = "";
            }
            if (message.username != null && message.hasOwnProperty("username"))
                object.username = message.username;
            if (message.token != null && message.hasOwnProperty("token"))
                object.token = message.token;
            if (message.uid != null && message.hasOwnProperty("uid"))
                object.uid = message.uid;
            let keys2;
            if (message.items && (keys2 = Object.keys(message.items)).length) {
                object.items = {};
                for (let j = 0; j < keys2.length; ++j)
                    if (typeof message.items[keys2[j]] === "number")
                        object.items[keys2[j]] = options.longs === String ? String(message.items[keys2[j]]) : message.items[keys2[j]];
                    else
                        object.items[keys2[j]] = options.longs === String ? $util.Long.prototype.toString.call(message.items[keys2[j]]) : options.longs === Number ? new $util.LongBits(message.items[keys2[j]].low >>> 0, message.items[keys2[j]].high >>> 0).toNumber() : message.items[keys2[j]];
            }
            return object;
        };

        /**
         * Converts this AccountLoginResp to JSON.
         * @function toJSON
         * @memberof ws.AccountLoginResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        AccountLoginResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return AccountLoginResp;
    })();

    ws.AccountRegister = (function() {

        /**
         * Properties of an AccountRegister.
         * @memberof ws
         * @interface IAccountRegister
         * @property {string|null} [username] AccountRegister username
         * @property {string|null} [password] AccountRegister password
         */

        /**
         * Constructs a new AccountRegister.
         * @memberof ws
         * @classdesc Represents an AccountRegister.
         * @implements IAccountRegister
         * @constructor
         * @param {ws.IAccountRegister=} [properties] Properties to set
         */
        function AccountRegister(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * AccountRegister username.
         * @member {string} username
         * @memberof ws.AccountRegister
         * @instance
         */
        AccountRegister.prototype.username = "";

        /**
         * AccountRegister password.
         * @member {string} password
         * @memberof ws.AccountRegister
         * @instance
         */
        AccountRegister.prototype.password = "";

        /**
         * Creates a new AccountRegister instance using the specified properties.
         * @function create
         * @memberof ws.AccountRegister
         * @static
         * @param {ws.IAccountRegister=} [properties] Properties to set
         * @returns {ws.AccountRegister} AccountRegister instance
         */
        AccountRegister.create = function create(properties) {
            return new AccountRegister(properties);
        };

        /**
         * Encodes the specified AccountRegister message. Does not implicitly {@link ws.AccountRegister.verify|verify} messages.
         * @function encode
         * @memberof ws.AccountRegister
         * @static
         * @param {ws.IAccountRegister} message AccountRegister message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountRegister.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.username != null && Object.hasOwnProperty.call(message, "username"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.username);
            if (message.password != null && Object.hasOwnProperty.call(message, "password"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.password);
            return writer;
        };

        /**
         * Encodes the specified AccountRegister message, length delimited. Does not implicitly {@link ws.AccountRegister.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.AccountRegister
         * @static
         * @param {ws.IAccountRegister} message AccountRegister message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountRegister.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an AccountRegister message from the specified reader or buffer.
         * @function decode
         * @memberof ws.AccountRegister
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.AccountRegister} AccountRegister
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountRegister.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.AccountRegister();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.username = reader.string();
                    break;
                case 2:
                    message.password = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an AccountRegister message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.AccountRegister
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.AccountRegister} AccountRegister
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountRegister.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an AccountRegister message.
         * @function verify
         * @memberof ws.AccountRegister
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        AccountRegister.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.username != null && message.hasOwnProperty("username"))
                if (!$util.isString(message.username))
                    return "username: string expected";
            if (message.password != null && message.hasOwnProperty("password"))
                if (!$util.isString(message.password))
                    return "password: string expected";
            return null;
        };

        /**
         * Creates an AccountRegister message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.AccountRegister
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.AccountRegister} AccountRegister
         */
        AccountRegister.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.AccountRegister)
                return object;
            let message = new $root.ws.AccountRegister();
            if (object.username != null)
                message.username = String(object.username);
            if (object.password != null)
                message.password = String(object.password);
            return message;
        };

        /**
         * Creates a plain object from an AccountRegister message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.AccountRegister
         * @static
         * @param {ws.AccountRegister} message AccountRegister
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        AccountRegister.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.username = "";
                object.password = "";
            }
            if (message.username != null && message.hasOwnProperty("username"))
                object.username = message.username;
            if (message.password != null && message.hasOwnProperty("password"))
                object.password = message.password;
            return object;
        };

        /**
         * Converts this AccountRegister to JSON.
         * @function toJSON
         * @memberof ws.AccountRegister
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        AccountRegister.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return AccountRegister;
    })();

    ws.AccountRegisterResp = (function() {

        /**
         * Properties of an AccountRegisterResp.
         * @memberof ws
         * @interface IAccountRegisterResp
         * @property {string|null} [username] AccountRegisterResp username
         * @property {string|null} [uid] AccountRegisterResp uid
         */

        /**
         * Constructs a new AccountRegisterResp.
         * @memberof ws
         * @classdesc Represents an AccountRegisterResp.
         * @implements IAccountRegisterResp
         * @constructor
         * @param {ws.IAccountRegisterResp=} [properties] Properties to set
         */
        function AccountRegisterResp(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * AccountRegisterResp username.
         * @member {string} username
         * @memberof ws.AccountRegisterResp
         * @instance
         */
        AccountRegisterResp.prototype.username = "";

        /**
         * AccountRegisterResp uid.
         * @member {string} uid
         * @memberof ws.AccountRegisterResp
         * @instance
         */
        AccountRegisterResp.prototype.uid = "";

        /**
         * Creates a new AccountRegisterResp instance using the specified properties.
         * @function create
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {ws.IAccountRegisterResp=} [properties] Properties to set
         * @returns {ws.AccountRegisterResp} AccountRegisterResp instance
         */
        AccountRegisterResp.create = function create(properties) {
            return new AccountRegisterResp(properties);
        };

        /**
         * Encodes the specified AccountRegisterResp message. Does not implicitly {@link ws.AccountRegisterResp.verify|verify} messages.
         * @function encode
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {ws.IAccountRegisterResp} message AccountRegisterResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountRegisterResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.username != null && Object.hasOwnProperty.call(message, "username"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.username);
            if (message.uid != null && Object.hasOwnProperty.call(message, "uid"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.uid);
            return writer;
        };

        /**
         * Encodes the specified AccountRegisterResp message, length delimited. Does not implicitly {@link ws.AccountRegisterResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {ws.IAccountRegisterResp} message AccountRegisterResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        AccountRegisterResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an AccountRegisterResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.AccountRegisterResp} AccountRegisterResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountRegisterResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.AccountRegisterResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.username = reader.string();
                    break;
                case 2:
                    message.uid = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an AccountRegisterResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.AccountRegisterResp} AccountRegisterResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        AccountRegisterResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an AccountRegisterResp message.
         * @function verify
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        AccountRegisterResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.username != null && message.hasOwnProperty("username"))
                if (!$util.isString(message.username))
                    return "username: string expected";
            if (message.uid != null && message.hasOwnProperty("uid"))
                if (!$util.isString(message.uid))
                    return "uid: string expected";
            return null;
        };

        /**
         * Creates an AccountRegisterResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.AccountRegisterResp} AccountRegisterResp
         */
        AccountRegisterResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.AccountRegisterResp)
                return object;
            let message = new $root.ws.AccountRegisterResp();
            if (object.username != null)
                message.username = String(object.username);
            if (object.uid != null)
                message.uid = String(object.uid);
            return message;
        };

        /**
         * Creates a plain object from an AccountRegisterResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.AccountRegisterResp
         * @static
         * @param {ws.AccountRegisterResp} message AccountRegisterResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        AccountRegisterResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.username = "";
                object.uid = "";
            }
            if (message.username != null && message.hasOwnProperty("username"))
                object.username = message.username;
            if (message.uid != null && message.hasOwnProperty("uid"))
                object.uid = message.uid;
            return object;
        };

        /**
         * Converts this AccountRegisterResp to JSON.
         * @function toJSON
         * @memberof ws.AccountRegisterResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        AccountRegisterResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return AccountRegisterResp;
    })();

    ws.ChatWorldSend = (function() {

        /**
         * Properties of a ChatWorldSend.
         * @memberof ws
         * @interface IChatWorldSend
         * @property {string|null} [content] ChatWorldSend content
         */

        /**
         * Constructs a new ChatWorldSend.
         * @memberof ws
         * @classdesc Represents a ChatWorldSend.
         * @implements IChatWorldSend
         * @constructor
         * @param {ws.IChatWorldSend=} [properties] Properties to set
         */
        function ChatWorldSend(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatWorldSend content.
         * @member {string} content
         * @memberof ws.ChatWorldSend
         * @instance
         */
        ChatWorldSend.prototype.content = "";

        /**
         * Creates a new ChatWorldSend instance using the specified properties.
         * @function create
         * @memberof ws.ChatWorldSend
         * @static
         * @param {ws.IChatWorldSend=} [properties] Properties to set
         * @returns {ws.ChatWorldSend} ChatWorldSend instance
         */
        ChatWorldSend.create = function create(properties) {
            return new ChatWorldSend(properties);
        };

        /**
         * Encodes the specified ChatWorldSend message. Does not implicitly {@link ws.ChatWorldSend.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatWorldSend
         * @static
         * @param {ws.IChatWorldSend} message ChatWorldSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatWorldSend.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.content);
            return writer;
        };

        /**
         * Encodes the specified ChatWorldSend message, length delimited. Does not implicitly {@link ws.ChatWorldSend.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatWorldSend
         * @static
         * @param {ws.IChatWorldSend} message ChatWorldSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatWorldSend.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatWorldSend message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatWorldSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatWorldSend} ChatWorldSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatWorldSend.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatWorldSend();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.content = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatWorldSend message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatWorldSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatWorldSend} ChatWorldSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatWorldSend.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatWorldSend message.
         * @function verify
         * @memberof ws.ChatWorldSend
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatWorldSend.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            return null;
        };

        /**
         * Creates a ChatWorldSend message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatWorldSend
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatWorldSend} ChatWorldSend
         */
        ChatWorldSend.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatWorldSend)
                return object;
            let message = new $root.ws.ChatWorldSend();
            if (object.content != null)
                message.content = String(object.content);
            return message;
        };

        /**
         * Creates a plain object from a ChatWorldSend message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatWorldSend
         * @static
         * @param {ws.ChatWorldSend} message ChatWorldSend
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatWorldSend.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.content = "";
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            return object;
        };

        /**
         * Converts this ChatWorldSend to JSON.
         * @function toJSON
         * @memberof ws.ChatWorldSend
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatWorldSend.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatWorldSend;
    })();

    ws.ChatPrivateSend = (function() {

        /**
         * Properties of a ChatPrivateSend.
         * @memberof ws
         * @interface IChatPrivateSend
         * @property {string|null} [toUid] ChatPrivateSend toUid
         * @property {string|null} [content] ChatPrivateSend content
         */

        /**
         * Constructs a new ChatPrivateSend.
         * @memberof ws
         * @classdesc Represents a ChatPrivateSend.
         * @implements IChatPrivateSend
         * @constructor
         * @param {ws.IChatPrivateSend=} [properties] Properties to set
         */
        function ChatPrivateSend(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatPrivateSend toUid.
         * @member {string} toUid
         * @memberof ws.ChatPrivateSend
         * @instance
         */
        ChatPrivateSend.prototype.toUid = "";

        /**
         * ChatPrivateSend content.
         * @member {string} content
         * @memberof ws.ChatPrivateSend
         * @instance
         */
        ChatPrivateSend.prototype.content = "";

        /**
         * Creates a new ChatPrivateSend instance using the specified properties.
         * @function create
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {ws.IChatPrivateSend=} [properties] Properties to set
         * @returns {ws.ChatPrivateSend} ChatPrivateSend instance
         */
        ChatPrivateSend.create = function create(properties) {
            return new ChatPrivateSend(properties);
        };

        /**
         * Encodes the specified ChatPrivateSend message. Does not implicitly {@link ws.ChatPrivateSend.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {ws.IChatPrivateSend} message ChatPrivateSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatPrivateSend.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.toUid != null && Object.hasOwnProperty.call(message, "toUid"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.toUid);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.content);
            return writer;
        };

        /**
         * Encodes the specified ChatPrivateSend message, length delimited. Does not implicitly {@link ws.ChatPrivateSend.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {ws.IChatPrivateSend} message ChatPrivateSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatPrivateSend.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatPrivateSend message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatPrivateSend} ChatPrivateSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatPrivateSend.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatPrivateSend();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.toUid = reader.string();
                    break;
                case 2:
                    message.content = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatPrivateSend message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatPrivateSend} ChatPrivateSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatPrivateSend.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatPrivateSend message.
         * @function verify
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatPrivateSend.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.toUid != null && message.hasOwnProperty("toUid"))
                if (!$util.isString(message.toUid))
                    return "toUid: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            return null;
        };

        /**
         * Creates a ChatPrivateSend message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatPrivateSend} ChatPrivateSend
         */
        ChatPrivateSend.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatPrivateSend)
                return object;
            let message = new $root.ws.ChatPrivateSend();
            if (object.toUid != null)
                message.toUid = String(object.toUid);
            if (object.content != null)
                message.content = String(object.content);
            return message;
        };

        /**
         * Creates a plain object from a ChatPrivateSend message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatPrivateSend
         * @static
         * @param {ws.ChatPrivateSend} message ChatPrivateSend
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatPrivateSend.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.toUid = "";
                object.content = "";
            }
            if (message.toUid != null && message.hasOwnProperty("toUid"))
                object.toUid = message.toUid;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            return object;
        };

        /**
         * Converts this ChatPrivateSend to JSON.
         * @function toJSON
         * @memberof ws.ChatPrivateSend
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatPrivateSend.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatPrivateSend;
    })();

    ws.ChatGroupSend = (function() {

        /**
         * Properties of a ChatGroupSend.
         * @memberof ws
         * @interface IChatGroupSend
         * @property {string|null} [groupId] ChatGroupSend groupId
         * @property {string|null} [content] ChatGroupSend content
         */

        /**
         * Constructs a new ChatGroupSend.
         * @memberof ws
         * @classdesc Represents a ChatGroupSend.
         * @implements IChatGroupSend
         * @constructor
         * @param {ws.IChatGroupSend=} [properties] Properties to set
         */
        function ChatGroupSend(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatGroupSend groupId.
         * @member {string} groupId
         * @memberof ws.ChatGroupSend
         * @instance
         */
        ChatGroupSend.prototype.groupId = "";

        /**
         * ChatGroupSend content.
         * @member {string} content
         * @memberof ws.ChatGroupSend
         * @instance
         */
        ChatGroupSend.prototype.content = "";

        /**
         * Creates a new ChatGroupSend instance using the specified properties.
         * @function create
         * @memberof ws.ChatGroupSend
         * @static
         * @param {ws.IChatGroupSend=} [properties] Properties to set
         * @returns {ws.ChatGroupSend} ChatGroupSend instance
         */
        ChatGroupSend.create = function create(properties) {
            return new ChatGroupSend(properties);
        };

        /**
         * Encodes the specified ChatGroupSend message. Does not implicitly {@link ws.ChatGroupSend.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatGroupSend
         * @static
         * @param {ws.IChatGroupSend} message ChatGroupSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatGroupSend.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.groupId != null && Object.hasOwnProperty.call(message, "groupId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.groupId);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.content);
            return writer;
        };

        /**
         * Encodes the specified ChatGroupSend message, length delimited. Does not implicitly {@link ws.ChatGroupSend.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatGroupSend
         * @static
         * @param {ws.IChatGroupSend} message ChatGroupSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatGroupSend.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatGroupSend message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatGroupSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatGroupSend} ChatGroupSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatGroupSend.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatGroupSend();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.groupId = reader.string();
                    break;
                case 2:
                    message.content = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatGroupSend message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatGroupSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatGroupSend} ChatGroupSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatGroupSend.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatGroupSend message.
         * @function verify
         * @memberof ws.ChatGroupSend
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatGroupSend.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                if (!$util.isString(message.groupId))
                    return "groupId: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            return null;
        };

        /**
         * Creates a ChatGroupSend message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatGroupSend
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatGroupSend} ChatGroupSend
         */
        ChatGroupSend.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatGroupSend)
                return object;
            let message = new $root.ws.ChatGroupSend();
            if (object.groupId != null)
                message.groupId = String(object.groupId);
            if (object.content != null)
                message.content = String(object.content);
            return message;
        };

        /**
         * Creates a plain object from a ChatGroupSend message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatGroupSend
         * @static
         * @param {ws.ChatGroupSend} message ChatGroupSend
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatGroupSend.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.groupId = "";
                object.content = "";
            }
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                object.groupId = message.groupId;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            return object;
        };

        /**
         * Converts this ChatGroupSend to JSON.
         * @function toJSON
         * @memberof ws.ChatGroupSend
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatGroupSend.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatGroupSend;
    })();

    ws.GroupMember = (function() {

        /**
         * Properties of a GroupMember.
         * @memberof ws
         * @interface IGroupMember
         * @property {string|null} [groupId] GroupMember groupId
         */

        /**
         * Constructs a new GroupMember.
         * @memberof ws
         * @classdesc Represents a GroupMember.
         * @implements IGroupMember
         * @constructor
         * @param {ws.IGroupMember=} [properties] Properties to set
         */
        function GroupMember(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * GroupMember groupId.
         * @member {string} groupId
         * @memberof ws.GroupMember
         * @instance
         */
        GroupMember.prototype.groupId = "";

        /**
         * Creates a new GroupMember instance using the specified properties.
         * @function create
         * @memberof ws.GroupMember
         * @static
         * @param {ws.IGroupMember=} [properties] Properties to set
         * @returns {ws.GroupMember} GroupMember instance
         */
        GroupMember.create = function create(properties) {
            return new GroupMember(properties);
        };

        /**
         * Encodes the specified GroupMember message. Does not implicitly {@link ws.GroupMember.verify|verify} messages.
         * @function encode
         * @memberof ws.GroupMember
         * @static
         * @param {ws.IGroupMember} message GroupMember message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GroupMember.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.groupId != null && Object.hasOwnProperty.call(message, "groupId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.groupId);
            return writer;
        };

        /**
         * Encodes the specified GroupMember message, length delimited. Does not implicitly {@link ws.GroupMember.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.GroupMember
         * @static
         * @param {ws.IGroupMember} message GroupMember message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GroupMember.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a GroupMember message from the specified reader or buffer.
         * @function decode
         * @memberof ws.GroupMember
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.GroupMember} GroupMember
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GroupMember.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.GroupMember();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.groupId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a GroupMember message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.GroupMember
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.GroupMember} GroupMember
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GroupMember.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a GroupMember message.
         * @function verify
         * @memberof ws.GroupMember
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GroupMember.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                if (!$util.isString(message.groupId))
                    return "groupId: string expected";
            return null;
        };

        /**
         * Creates a GroupMember message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.GroupMember
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.GroupMember} GroupMember
         */
        GroupMember.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.GroupMember)
                return object;
            let message = new $root.ws.GroupMember();
            if (object.groupId != null)
                message.groupId = String(object.groupId);
            return message;
        };

        /**
         * Creates a plain object from a GroupMember message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.GroupMember
         * @static
         * @param {ws.GroupMember} message GroupMember
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GroupMember.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.groupId = "";
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                object.groupId = message.groupId;
            return object;
        };

        /**
         * Converts this GroupMember to JSON.
         * @function toJSON
         * @memberof ws.GroupMember
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GroupMember.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return GroupMember;
    })();

    ws.GroupMemberResp = (function() {

        /**
         * Properties of a GroupMemberResp.
         * @memberof ws
         * @interface IGroupMemberResp
         * @property {string|null} [groupId] GroupMemberResp groupId
         * @property {Array.<string>|null} [members] GroupMemberResp members
         * @property {number|null} [count] GroupMemberResp count
         */

        /**
         * Constructs a new GroupMemberResp.
         * @memberof ws
         * @classdesc Represents a GroupMemberResp.
         * @implements IGroupMemberResp
         * @constructor
         * @param {ws.IGroupMemberResp=} [properties] Properties to set
         */
        function GroupMemberResp(properties) {
            this.members = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * GroupMemberResp groupId.
         * @member {string} groupId
         * @memberof ws.GroupMemberResp
         * @instance
         */
        GroupMemberResp.prototype.groupId = "";

        /**
         * GroupMemberResp members.
         * @member {Array.<string>} members
         * @memberof ws.GroupMemberResp
         * @instance
         */
        GroupMemberResp.prototype.members = $util.emptyArray;

        /**
         * GroupMemberResp count.
         * @member {number} count
         * @memberof ws.GroupMemberResp
         * @instance
         */
        GroupMemberResp.prototype.count = 0;

        /**
         * Creates a new GroupMemberResp instance using the specified properties.
         * @function create
         * @memberof ws.GroupMemberResp
         * @static
         * @param {ws.IGroupMemberResp=} [properties] Properties to set
         * @returns {ws.GroupMemberResp} GroupMemberResp instance
         */
        GroupMemberResp.create = function create(properties) {
            return new GroupMemberResp(properties);
        };

        /**
         * Encodes the specified GroupMemberResp message. Does not implicitly {@link ws.GroupMemberResp.verify|verify} messages.
         * @function encode
         * @memberof ws.GroupMemberResp
         * @static
         * @param {ws.IGroupMemberResp} message GroupMemberResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GroupMemberResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.groupId != null && Object.hasOwnProperty.call(message, "groupId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.groupId);
            if (message.members != null && message.members.length)
                for (let i = 0; i < message.members.length; ++i)
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message.members[i]);
            if (message.count != null && Object.hasOwnProperty.call(message, "count"))
                writer.uint32(/* id 3, wireType 0 =*/24).int32(message.count);
            return writer;
        };

        /**
         * Encodes the specified GroupMemberResp message, length delimited. Does not implicitly {@link ws.GroupMemberResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.GroupMemberResp
         * @static
         * @param {ws.IGroupMemberResp} message GroupMemberResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        GroupMemberResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a GroupMemberResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.GroupMemberResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.GroupMemberResp} GroupMemberResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GroupMemberResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.GroupMemberResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.groupId = reader.string();
                    break;
                case 2:
                    if (!(message.members && message.members.length))
                        message.members = [];
                    message.members.push(reader.string());
                    break;
                case 3:
                    message.count = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a GroupMemberResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.GroupMemberResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.GroupMemberResp} GroupMemberResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        GroupMemberResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a GroupMemberResp message.
         * @function verify
         * @memberof ws.GroupMemberResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        GroupMemberResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                if (!$util.isString(message.groupId))
                    return "groupId: string expected";
            if (message.members != null && message.hasOwnProperty("members")) {
                if (!Array.isArray(message.members))
                    return "members: array expected";
                for (let i = 0; i < message.members.length; ++i)
                    if (!$util.isString(message.members[i]))
                        return "members: string[] expected";
            }
            if (message.count != null && message.hasOwnProperty("count"))
                if (!$util.isInteger(message.count))
                    return "count: integer expected";
            return null;
        };

        /**
         * Creates a GroupMemberResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.GroupMemberResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.GroupMemberResp} GroupMemberResp
         */
        GroupMemberResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.GroupMemberResp)
                return object;
            let message = new $root.ws.GroupMemberResp();
            if (object.groupId != null)
                message.groupId = String(object.groupId);
            if (object.members) {
                if (!Array.isArray(object.members))
                    throw TypeError(".ws.GroupMemberResp.members: array expected");
                message.members = [];
                for (let i = 0; i < object.members.length; ++i)
                    message.members[i] = String(object.members[i]);
            }
            if (object.count != null)
                message.count = object.count | 0;
            return message;
        };

        /**
         * Creates a plain object from a GroupMemberResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.GroupMemberResp
         * @static
         * @param {ws.GroupMemberResp} message GroupMemberResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        GroupMemberResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.members = [];
            if (options.defaults) {
                object.groupId = "";
                object.count = 0;
            }
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                object.groupId = message.groupId;
            if (message.members && message.members.length) {
                object.members = [];
                for (let j = 0; j < message.members.length; ++j)
                    object.members[j] = message.members[j];
            }
            if (message.count != null && message.hasOwnProperty("count"))
                object.count = message.count;
            return object;
        };

        /**
         * Converts this GroupMemberResp to JSON.
         * @function toJSON
         * @memberof ws.GroupMemberResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        GroupMemberResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return GroupMemberResp;
    })();

    ws.ChatOnlineResp = (function() {

        /**
         * Properties of a ChatOnlineResp.
         * @memberof ws
         * @interface IChatOnlineResp
         * @property {number|null} [total] ChatOnlineResp total
         * @property {Array.<string>|null} [uids] ChatOnlineResp uids
         */

        /**
         * Constructs a new ChatOnlineResp.
         * @memberof ws
         * @classdesc Represents a ChatOnlineResp.
         * @implements IChatOnlineResp
         * @constructor
         * @param {ws.IChatOnlineResp=} [properties] Properties to set
         */
        function ChatOnlineResp(properties) {
            this.uids = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatOnlineResp total.
         * @member {number} total
         * @memberof ws.ChatOnlineResp
         * @instance
         */
        ChatOnlineResp.prototype.total = 0;

        /**
         * ChatOnlineResp uids.
         * @member {Array.<string>} uids
         * @memberof ws.ChatOnlineResp
         * @instance
         */
        ChatOnlineResp.prototype.uids = $util.emptyArray;

        /**
         * Creates a new ChatOnlineResp instance using the specified properties.
         * @function create
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {ws.IChatOnlineResp=} [properties] Properties to set
         * @returns {ws.ChatOnlineResp} ChatOnlineResp instance
         */
        ChatOnlineResp.create = function create(properties) {
            return new ChatOnlineResp(properties);
        };

        /**
         * Encodes the specified ChatOnlineResp message. Does not implicitly {@link ws.ChatOnlineResp.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {ws.IChatOnlineResp} message ChatOnlineResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatOnlineResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.total != null && Object.hasOwnProperty.call(message, "total"))
                writer.uint32(/* id 1, wireType 0 =*/8).int32(message.total);
            if (message.uids != null && message.uids.length)
                for (let i = 0; i < message.uids.length; ++i)
                    writer.uint32(/* id 2, wireType 2 =*/18).string(message.uids[i]);
            return writer;
        };

        /**
         * Encodes the specified ChatOnlineResp message, length delimited. Does not implicitly {@link ws.ChatOnlineResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {ws.IChatOnlineResp} message ChatOnlineResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatOnlineResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatOnlineResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatOnlineResp} ChatOnlineResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatOnlineResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatOnlineResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.total = reader.int32();
                    break;
                case 2:
                    if (!(message.uids && message.uids.length))
                        message.uids = [];
                    message.uids.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatOnlineResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatOnlineResp} ChatOnlineResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatOnlineResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatOnlineResp message.
         * @function verify
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatOnlineResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.total != null && message.hasOwnProperty("total"))
                if (!$util.isInteger(message.total))
                    return "total: integer expected";
            if (message.uids != null && message.hasOwnProperty("uids")) {
                if (!Array.isArray(message.uids))
                    return "uids: array expected";
                for (let i = 0; i < message.uids.length; ++i)
                    if (!$util.isString(message.uids[i]))
                        return "uids: string[] expected";
            }
            return null;
        };

        /**
         * Creates a ChatOnlineResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatOnlineResp} ChatOnlineResp
         */
        ChatOnlineResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatOnlineResp)
                return object;
            let message = new $root.ws.ChatOnlineResp();
            if (object.total != null)
                message.total = object.total | 0;
            if (object.uids) {
                if (!Array.isArray(object.uids))
                    throw TypeError(".ws.ChatOnlineResp.uids: array expected");
                message.uids = [];
                for (let i = 0; i < object.uids.length; ++i)
                    message.uids[i] = String(object.uids[i]);
            }
            return message;
        };

        /**
         * Creates a plain object from a ChatOnlineResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatOnlineResp
         * @static
         * @param {ws.ChatOnlineResp} message ChatOnlineResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatOnlineResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.uids = [];
            if (options.defaults)
                object.total = 0;
            if (message.total != null && message.hasOwnProperty("total"))
                object.total = message.total;
            if (message.uids && message.uids.length) {
                object.uids = [];
                for (let j = 0; j < message.uids.length; ++j)
                    object.uids[j] = message.uids[j];
            }
            return object;
        };

        /**
         * Converts this ChatOnlineResp to JSON.
         * @function toJSON
         * @memberof ws.ChatOnlineResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatOnlineResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatOnlineResp;
    })();

    ws.ChatWorldNew = (function() {

        /**
         * Properties of a ChatWorldNew.
         * @memberof ws
         * @interface IChatWorldNew
         * @property {string|null} [msgId] ChatWorldNew msgId
         * @property {string|null} [fromUid] ChatWorldNew fromUid
         * @property {string|null} [fromName] ChatWorldNew fromName
         * @property {string|null} [content] ChatWorldNew content
         * @property {number|Long|null} [timestamp] ChatWorldNew timestamp
         */

        /**
         * Constructs a new ChatWorldNew.
         * @memberof ws
         * @classdesc Represents a ChatWorldNew.
         * @implements IChatWorldNew
         * @constructor
         * @param {ws.IChatWorldNew=} [properties] Properties to set
         */
        function ChatWorldNew(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatWorldNew msgId.
         * @member {string} msgId
         * @memberof ws.ChatWorldNew
         * @instance
         */
        ChatWorldNew.prototype.msgId = "";

        /**
         * ChatWorldNew fromUid.
         * @member {string} fromUid
         * @memberof ws.ChatWorldNew
         * @instance
         */
        ChatWorldNew.prototype.fromUid = "";

        /**
         * ChatWorldNew fromName.
         * @member {string} fromName
         * @memberof ws.ChatWorldNew
         * @instance
         */
        ChatWorldNew.prototype.fromName = "";

        /**
         * ChatWorldNew content.
         * @member {string} content
         * @memberof ws.ChatWorldNew
         * @instance
         */
        ChatWorldNew.prototype.content = "";

        /**
         * ChatWorldNew timestamp.
         * @member {number|Long} timestamp
         * @memberof ws.ChatWorldNew
         * @instance
         */
        ChatWorldNew.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new ChatWorldNew instance using the specified properties.
         * @function create
         * @memberof ws.ChatWorldNew
         * @static
         * @param {ws.IChatWorldNew=} [properties] Properties to set
         * @returns {ws.ChatWorldNew} ChatWorldNew instance
         */
        ChatWorldNew.create = function create(properties) {
            return new ChatWorldNew(properties);
        };

        /**
         * Encodes the specified ChatWorldNew message. Does not implicitly {@link ws.ChatWorldNew.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatWorldNew
         * @static
         * @param {ws.IChatWorldNew} message ChatWorldNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatWorldNew.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.msgId != null && Object.hasOwnProperty.call(message, "msgId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.msgId);
            if (message.fromUid != null && Object.hasOwnProperty.call(message, "fromUid"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.fromUid);
            if (message.fromName != null && Object.hasOwnProperty.call(message, "fromName"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.fromName);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.content);
            if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
                writer.uint32(/* id 5, wireType 0 =*/40).int64(message.timestamp);
            return writer;
        };

        /**
         * Encodes the specified ChatWorldNew message, length delimited. Does not implicitly {@link ws.ChatWorldNew.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatWorldNew
         * @static
         * @param {ws.IChatWorldNew} message ChatWorldNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatWorldNew.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatWorldNew message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatWorldNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatWorldNew} ChatWorldNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatWorldNew.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatWorldNew();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.msgId = reader.string();
                    break;
                case 2:
                    message.fromUid = reader.string();
                    break;
                case 3:
                    message.fromName = reader.string();
                    break;
                case 4:
                    message.content = reader.string();
                    break;
                case 5:
                    message.timestamp = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatWorldNew message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatWorldNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatWorldNew} ChatWorldNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatWorldNew.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatWorldNew message.
         * @function verify
         * @memberof ws.ChatWorldNew
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatWorldNew.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                if (!$util.isString(message.msgId))
                    return "msgId: string expected";
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                if (!$util.isString(message.fromUid))
                    return "fromUid: string expected";
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                if (!$util.isString(message.fromName))
                    return "fromName: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                    return "timestamp: integer|Long expected";
            return null;
        };

        /**
         * Creates a ChatWorldNew message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatWorldNew
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatWorldNew} ChatWorldNew
         */
        ChatWorldNew.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatWorldNew)
                return object;
            let message = new $root.ws.ChatWorldNew();
            if (object.msgId != null)
                message.msgId = String(object.msgId);
            if (object.fromUid != null)
                message.fromUid = String(object.fromUid);
            if (object.fromName != null)
                message.fromName = String(object.fromName);
            if (object.content != null)
                message.content = String(object.content);
            if (object.timestamp != null)
                if ($util.Long)
                    (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
                else if (typeof object.timestamp === "string")
                    message.timestamp = parseInt(object.timestamp, 10);
                else if (typeof object.timestamp === "number")
                    message.timestamp = object.timestamp;
                else if (typeof object.timestamp === "object")
                    message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a ChatWorldNew message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatWorldNew
         * @static
         * @param {ws.ChatWorldNew} message ChatWorldNew
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatWorldNew.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.msgId = "";
                object.fromUid = "";
                object.fromName = "";
                object.content = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.timestamp = options.longs === String ? "0" : 0;
            }
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                object.msgId = message.msgId;
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                object.fromUid = message.fromUid;
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                object.fromName = message.fromName;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (typeof message.timestamp === "number")
                    object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
                else
                    object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
            return object;
        };

        /**
         * Converts this ChatWorldNew to JSON.
         * @function toJSON
         * @memberof ws.ChatWorldNew
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatWorldNew.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatWorldNew;
    })();

    ws.ChatPrivateNew = (function() {

        /**
         * Properties of a ChatPrivateNew.
         * @memberof ws
         * @interface IChatPrivateNew
         * @property {string|null} [msgId] ChatPrivateNew msgId
         * @property {string|null} [fromUid] ChatPrivateNew fromUid
         * @property {string|null} [fromName] ChatPrivateNew fromName
         * @property {string|null} [content] ChatPrivateNew content
         * @property {number|Long|null} [timestamp] ChatPrivateNew timestamp
         */

        /**
         * Constructs a new ChatPrivateNew.
         * @memberof ws
         * @classdesc Represents a ChatPrivateNew.
         * @implements IChatPrivateNew
         * @constructor
         * @param {ws.IChatPrivateNew=} [properties] Properties to set
         */
        function ChatPrivateNew(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatPrivateNew msgId.
         * @member {string} msgId
         * @memberof ws.ChatPrivateNew
         * @instance
         */
        ChatPrivateNew.prototype.msgId = "";

        /**
         * ChatPrivateNew fromUid.
         * @member {string} fromUid
         * @memberof ws.ChatPrivateNew
         * @instance
         */
        ChatPrivateNew.prototype.fromUid = "";

        /**
         * ChatPrivateNew fromName.
         * @member {string} fromName
         * @memberof ws.ChatPrivateNew
         * @instance
         */
        ChatPrivateNew.prototype.fromName = "";

        /**
         * ChatPrivateNew content.
         * @member {string} content
         * @memberof ws.ChatPrivateNew
         * @instance
         */
        ChatPrivateNew.prototype.content = "";

        /**
         * ChatPrivateNew timestamp.
         * @member {number|Long} timestamp
         * @memberof ws.ChatPrivateNew
         * @instance
         */
        ChatPrivateNew.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new ChatPrivateNew instance using the specified properties.
         * @function create
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {ws.IChatPrivateNew=} [properties] Properties to set
         * @returns {ws.ChatPrivateNew} ChatPrivateNew instance
         */
        ChatPrivateNew.create = function create(properties) {
            return new ChatPrivateNew(properties);
        };

        /**
         * Encodes the specified ChatPrivateNew message. Does not implicitly {@link ws.ChatPrivateNew.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {ws.IChatPrivateNew} message ChatPrivateNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatPrivateNew.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.msgId != null && Object.hasOwnProperty.call(message, "msgId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.msgId);
            if (message.fromUid != null && Object.hasOwnProperty.call(message, "fromUid"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.fromUid);
            if (message.fromName != null && Object.hasOwnProperty.call(message, "fromName"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.fromName);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.content);
            if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
                writer.uint32(/* id 5, wireType 0 =*/40).int64(message.timestamp);
            return writer;
        };

        /**
         * Encodes the specified ChatPrivateNew message, length delimited. Does not implicitly {@link ws.ChatPrivateNew.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {ws.IChatPrivateNew} message ChatPrivateNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatPrivateNew.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatPrivateNew message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatPrivateNew} ChatPrivateNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatPrivateNew.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatPrivateNew();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.msgId = reader.string();
                    break;
                case 2:
                    message.fromUid = reader.string();
                    break;
                case 3:
                    message.fromName = reader.string();
                    break;
                case 4:
                    message.content = reader.string();
                    break;
                case 5:
                    message.timestamp = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatPrivateNew message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatPrivateNew} ChatPrivateNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatPrivateNew.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatPrivateNew message.
         * @function verify
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatPrivateNew.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                if (!$util.isString(message.msgId))
                    return "msgId: string expected";
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                if (!$util.isString(message.fromUid))
                    return "fromUid: string expected";
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                if (!$util.isString(message.fromName))
                    return "fromName: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                    return "timestamp: integer|Long expected";
            return null;
        };

        /**
         * Creates a ChatPrivateNew message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatPrivateNew} ChatPrivateNew
         */
        ChatPrivateNew.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatPrivateNew)
                return object;
            let message = new $root.ws.ChatPrivateNew();
            if (object.msgId != null)
                message.msgId = String(object.msgId);
            if (object.fromUid != null)
                message.fromUid = String(object.fromUid);
            if (object.fromName != null)
                message.fromName = String(object.fromName);
            if (object.content != null)
                message.content = String(object.content);
            if (object.timestamp != null)
                if ($util.Long)
                    (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
                else if (typeof object.timestamp === "string")
                    message.timestamp = parseInt(object.timestamp, 10);
                else if (typeof object.timestamp === "number")
                    message.timestamp = object.timestamp;
                else if (typeof object.timestamp === "object")
                    message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a ChatPrivateNew message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatPrivateNew
         * @static
         * @param {ws.ChatPrivateNew} message ChatPrivateNew
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatPrivateNew.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.msgId = "";
                object.fromUid = "";
                object.fromName = "";
                object.content = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.timestamp = options.longs === String ? "0" : 0;
            }
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                object.msgId = message.msgId;
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                object.fromUid = message.fromUid;
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                object.fromName = message.fromName;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (typeof message.timestamp === "number")
                    object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
                else
                    object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
            return object;
        };

        /**
         * Converts this ChatPrivateNew to JSON.
         * @function toJSON
         * @memberof ws.ChatPrivateNew
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatPrivateNew.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatPrivateNew;
    })();

    ws.ChatGroupNew = (function() {

        /**
         * Properties of a ChatGroupNew.
         * @memberof ws
         * @interface IChatGroupNew
         * @property {string|null} [msgId] ChatGroupNew msgId
         * @property {string|null} [groupId] ChatGroupNew groupId
         * @property {string|null} [fromUid] ChatGroupNew fromUid
         * @property {string|null} [fromName] ChatGroupNew fromName
         * @property {string|null} [content] ChatGroupNew content
         * @property {number|Long|null} [timestamp] ChatGroupNew timestamp
         */

        /**
         * Constructs a new ChatGroupNew.
         * @memberof ws
         * @classdesc Represents a ChatGroupNew.
         * @implements IChatGroupNew
         * @constructor
         * @param {ws.IChatGroupNew=} [properties] Properties to set
         */
        function ChatGroupNew(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ChatGroupNew msgId.
         * @member {string} msgId
         * @memberof ws.ChatGroupNew
         * @instance
         */
        ChatGroupNew.prototype.msgId = "";

        /**
         * ChatGroupNew groupId.
         * @member {string} groupId
         * @memberof ws.ChatGroupNew
         * @instance
         */
        ChatGroupNew.prototype.groupId = "";

        /**
         * ChatGroupNew fromUid.
         * @member {string} fromUid
         * @memberof ws.ChatGroupNew
         * @instance
         */
        ChatGroupNew.prototype.fromUid = "";

        /**
         * ChatGroupNew fromName.
         * @member {string} fromName
         * @memberof ws.ChatGroupNew
         * @instance
         */
        ChatGroupNew.prototype.fromName = "";

        /**
         * ChatGroupNew content.
         * @member {string} content
         * @memberof ws.ChatGroupNew
         * @instance
         */
        ChatGroupNew.prototype.content = "";

        /**
         * ChatGroupNew timestamp.
         * @member {number|Long} timestamp
         * @memberof ws.ChatGroupNew
         * @instance
         */
        ChatGroupNew.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new ChatGroupNew instance using the specified properties.
         * @function create
         * @memberof ws.ChatGroupNew
         * @static
         * @param {ws.IChatGroupNew=} [properties] Properties to set
         * @returns {ws.ChatGroupNew} ChatGroupNew instance
         */
        ChatGroupNew.create = function create(properties) {
            return new ChatGroupNew(properties);
        };

        /**
         * Encodes the specified ChatGroupNew message. Does not implicitly {@link ws.ChatGroupNew.verify|verify} messages.
         * @function encode
         * @memberof ws.ChatGroupNew
         * @static
         * @param {ws.IChatGroupNew} message ChatGroupNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatGroupNew.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.msgId != null && Object.hasOwnProperty.call(message, "msgId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.msgId);
            if (message.groupId != null && Object.hasOwnProperty.call(message, "groupId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.groupId);
            if (message.fromUid != null && Object.hasOwnProperty.call(message, "fromUid"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.fromUid);
            if (message.fromName != null && Object.hasOwnProperty.call(message, "fromName"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.fromName);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.content);
            if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
                writer.uint32(/* id 6, wireType 0 =*/48).int64(message.timestamp);
            return writer;
        };

        /**
         * Encodes the specified ChatGroupNew message, length delimited. Does not implicitly {@link ws.ChatGroupNew.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ChatGroupNew
         * @static
         * @param {ws.IChatGroupNew} message ChatGroupNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ChatGroupNew.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a ChatGroupNew message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ChatGroupNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ChatGroupNew} ChatGroupNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatGroupNew.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ChatGroupNew();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.msgId = reader.string();
                    break;
                case 2:
                    message.groupId = reader.string();
                    break;
                case 3:
                    message.fromUid = reader.string();
                    break;
                case 4:
                    message.fromName = reader.string();
                    break;
                case 5:
                    message.content = reader.string();
                    break;
                case 6:
                    message.timestamp = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a ChatGroupNew message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ChatGroupNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ChatGroupNew} ChatGroupNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ChatGroupNew.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a ChatGroupNew message.
         * @function verify
         * @memberof ws.ChatGroupNew
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ChatGroupNew.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                if (!$util.isString(message.msgId))
                    return "msgId: string expected";
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                if (!$util.isString(message.groupId))
                    return "groupId: string expected";
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                if (!$util.isString(message.fromUid))
                    return "fromUid: string expected";
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                if (!$util.isString(message.fromName))
                    return "fromName: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                    return "timestamp: integer|Long expected";
            return null;
        };

        /**
         * Creates a ChatGroupNew message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ChatGroupNew
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ChatGroupNew} ChatGroupNew
         */
        ChatGroupNew.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ChatGroupNew)
                return object;
            let message = new $root.ws.ChatGroupNew();
            if (object.msgId != null)
                message.msgId = String(object.msgId);
            if (object.groupId != null)
                message.groupId = String(object.groupId);
            if (object.fromUid != null)
                message.fromUid = String(object.fromUid);
            if (object.fromName != null)
                message.fromName = String(object.fromName);
            if (object.content != null)
                message.content = String(object.content);
            if (object.timestamp != null)
                if ($util.Long)
                    (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
                else if (typeof object.timestamp === "string")
                    message.timestamp = parseInt(object.timestamp, 10);
                else if (typeof object.timestamp === "number")
                    message.timestamp = object.timestamp;
                else if (typeof object.timestamp === "object")
                    message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a ChatGroupNew message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ChatGroupNew
         * @static
         * @param {ws.ChatGroupNew} message ChatGroupNew
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ChatGroupNew.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.msgId = "";
                object.groupId = "";
                object.fromUid = "";
                object.fromName = "";
                object.content = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.timestamp = options.longs === String ? "0" : 0;
            }
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                object.msgId = message.msgId;
            if (message.groupId != null && message.hasOwnProperty("groupId"))
                object.groupId = message.groupId;
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                object.fromUid = message.fromUid;
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                object.fromName = message.fromName;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (typeof message.timestamp === "number")
                    object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
                else
                    object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
            return object;
        };

        /**
         * Converts this ChatGroupNew to JSON.
         * @function toJSON
         * @memberof ws.ChatGroupNew
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ChatGroupNew.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ChatGroupNew;
    })();

    ws.RoomCreate = (function() {

        /**
         * Properties of a RoomCreate.
         * @memberof ws
         * @interface IRoomCreate
         * @property {string|null} [name] RoomCreate name
         * @property {number|null} [maxSize] RoomCreate maxSize
         */

        /**
         * Constructs a new RoomCreate.
         * @memberof ws
         * @classdesc Represents a RoomCreate.
         * @implements IRoomCreate
         * @constructor
         * @param {ws.IRoomCreate=} [properties] Properties to set
         */
        function RoomCreate(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomCreate name.
         * @member {string} name
         * @memberof ws.RoomCreate
         * @instance
         */
        RoomCreate.prototype.name = "";

        /**
         * RoomCreate maxSize.
         * @member {number} maxSize
         * @memberof ws.RoomCreate
         * @instance
         */
        RoomCreate.prototype.maxSize = 0;

        /**
         * Creates a new RoomCreate instance using the specified properties.
         * @function create
         * @memberof ws.RoomCreate
         * @static
         * @param {ws.IRoomCreate=} [properties] Properties to set
         * @returns {ws.RoomCreate} RoomCreate instance
         */
        RoomCreate.create = function create(properties) {
            return new RoomCreate(properties);
        };

        /**
         * Encodes the specified RoomCreate message. Does not implicitly {@link ws.RoomCreate.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomCreate
         * @static
         * @param {ws.IRoomCreate} message RoomCreate message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomCreate.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.name);
            if (message.maxSize != null && Object.hasOwnProperty.call(message, "maxSize"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.maxSize);
            return writer;
        };

        /**
         * Encodes the specified RoomCreate message, length delimited. Does not implicitly {@link ws.RoomCreate.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomCreate
         * @static
         * @param {ws.IRoomCreate} message RoomCreate message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomCreate.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomCreate message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomCreate
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomCreate} RoomCreate
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomCreate.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomCreate();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.name = reader.string();
                    break;
                case 2:
                    message.maxSize = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomCreate message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomCreate
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomCreate} RoomCreate
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomCreate.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomCreate message.
         * @function verify
         * @memberof ws.RoomCreate
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomCreate.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            if (message.maxSize != null && message.hasOwnProperty("maxSize"))
                if (!$util.isInteger(message.maxSize))
                    return "maxSize: integer expected";
            return null;
        };

        /**
         * Creates a RoomCreate message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomCreate
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomCreate} RoomCreate
         */
        RoomCreate.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomCreate)
                return object;
            let message = new $root.ws.RoomCreate();
            if (object.name != null)
                message.name = String(object.name);
            if (object.maxSize != null)
                message.maxSize = object.maxSize | 0;
            return message;
        };

        /**
         * Creates a plain object from a RoomCreate message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomCreate
         * @static
         * @param {ws.RoomCreate} message RoomCreate
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomCreate.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.name = "";
                object.maxSize = 0;
            }
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            if (message.maxSize != null && message.hasOwnProperty("maxSize"))
                object.maxSize = message.maxSize;
            return object;
        };

        /**
         * Converts this RoomCreate to JSON.
         * @function toJSON
         * @memberof ws.RoomCreate
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomCreate.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomCreate;
    })();

    ws.RoomCreateResp = (function() {

        /**
         * Properties of a RoomCreateResp.
         * @memberof ws
         * @interface IRoomCreateResp
         * @property {string|null} [roomId] RoomCreateResp roomId
         * @property {string|null} [owner] RoomCreateResp owner
         * @property {string|null} [name] RoomCreateResp name
         */

        /**
         * Constructs a new RoomCreateResp.
         * @memberof ws
         * @classdesc Represents a RoomCreateResp.
         * @implements IRoomCreateResp
         * @constructor
         * @param {ws.IRoomCreateResp=} [properties] Properties to set
         */
        function RoomCreateResp(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomCreateResp roomId.
         * @member {string} roomId
         * @memberof ws.RoomCreateResp
         * @instance
         */
        RoomCreateResp.prototype.roomId = "";

        /**
         * RoomCreateResp owner.
         * @member {string} owner
         * @memberof ws.RoomCreateResp
         * @instance
         */
        RoomCreateResp.prototype.owner = "";

        /**
         * RoomCreateResp name.
         * @member {string} name
         * @memberof ws.RoomCreateResp
         * @instance
         */
        RoomCreateResp.prototype.name = "";

        /**
         * Creates a new RoomCreateResp instance using the specified properties.
         * @function create
         * @memberof ws.RoomCreateResp
         * @static
         * @param {ws.IRoomCreateResp=} [properties] Properties to set
         * @returns {ws.RoomCreateResp} RoomCreateResp instance
         */
        RoomCreateResp.create = function create(properties) {
            return new RoomCreateResp(properties);
        };

        /**
         * Encodes the specified RoomCreateResp message. Does not implicitly {@link ws.RoomCreateResp.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomCreateResp
         * @static
         * @param {ws.IRoomCreateResp} message RoomCreateResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomCreateResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            if (message.owner != null && Object.hasOwnProperty.call(message, "owner"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.owner);
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.name);
            return writer;
        };

        /**
         * Encodes the specified RoomCreateResp message, length delimited. Does not implicitly {@link ws.RoomCreateResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomCreateResp
         * @static
         * @param {ws.IRoomCreateResp} message RoomCreateResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomCreateResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomCreateResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomCreateResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomCreateResp} RoomCreateResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomCreateResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomCreateResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                case 2:
                    message.owner = reader.string();
                    break;
                case 3:
                    message.name = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomCreateResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomCreateResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomCreateResp} RoomCreateResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomCreateResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomCreateResp message.
         * @function verify
         * @memberof ws.RoomCreateResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomCreateResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            if (message.owner != null && message.hasOwnProperty("owner"))
                if (!$util.isString(message.owner))
                    return "owner: string expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            return null;
        };

        /**
         * Creates a RoomCreateResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomCreateResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomCreateResp} RoomCreateResp
         */
        RoomCreateResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomCreateResp)
                return object;
            let message = new $root.ws.RoomCreateResp();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            if (object.owner != null)
                message.owner = String(object.owner);
            if (object.name != null)
                message.name = String(object.name);
            return message;
        };

        /**
         * Creates a plain object from a RoomCreateResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomCreateResp
         * @static
         * @param {ws.RoomCreateResp} message RoomCreateResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomCreateResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.roomId = "";
                object.owner = "";
                object.name = "";
            }
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            if (message.owner != null && message.hasOwnProperty("owner"))
                object.owner = message.owner;
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            return object;
        };

        /**
         * Converts this RoomCreateResp to JSON.
         * @function toJSON
         * @memberof ws.RoomCreateResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomCreateResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomCreateResp;
    })();

    ws.RoomJoin = (function() {

        /**
         * Properties of a RoomJoin.
         * @memberof ws
         * @interface IRoomJoin
         * @property {string|null} [roomId] RoomJoin roomId
         */

        /**
         * Constructs a new RoomJoin.
         * @memberof ws
         * @classdesc Represents a RoomJoin.
         * @implements IRoomJoin
         * @constructor
         * @param {ws.IRoomJoin=} [properties] Properties to set
         */
        function RoomJoin(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomJoin roomId.
         * @member {string} roomId
         * @memberof ws.RoomJoin
         * @instance
         */
        RoomJoin.prototype.roomId = "";

        /**
         * Creates a new RoomJoin instance using the specified properties.
         * @function create
         * @memberof ws.RoomJoin
         * @static
         * @param {ws.IRoomJoin=} [properties] Properties to set
         * @returns {ws.RoomJoin} RoomJoin instance
         */
        RoomJoin.create = function create(properties) {
            return new RoomJoin(properties);
        };

        /**
         * Encodes the specified RoomJoin message. Does not implicitly {@link ws.RoomJoin.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomJoin
         * @static
         * @param {ws.IRoomJoin} message RoomJoin message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomJoin.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            return writer;
        };

        /**
         * Encodes the specified RoomJoin message, length delimited. Does not implicitly {@link ws.RoomJoin.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomJoin
         * @static
         * @param {ws.IRoomJoin} message RoomJoin message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomJoin.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomJoin message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomJoin
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomJoin} RoomJoin
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomJoin.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomJoin();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomJoin message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomJoin
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomJoin} RoomJoin
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomJoin.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomJoin message.
         * @function verify
         * @memberof ws.RoomJoin
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomJoin.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            return null;
        };

        /**
         * Creates a RoomJoin message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomJoin
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomJoin} RoomJoin
         */
        RoomJoin.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomJoin)
                return object;
            let message = new $root.ws.RoomJoin();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            return message;
        };

        /**
         * Creates a plain object from a RoomJoin message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomJoin
         * @static
         * @param {ws.RoomJoin} message RoomJoin
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomJoin.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.roomId = "";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            return object;
        };

        /**
         * Converts this RoomJoin to JSON.
         * @function toJSON
         * @memberof ws.RoomJoin
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomJoin.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomJoin;
    })();

    ws.RoomJoinResp = (function() {

        /**
         * Properties of a RoomJoinResp.
         * @memberof ws
         * @interface IRoomJoinResp
         * @property {string|null} [roomId] RoomJoinResp roomId
         * @property {string|null} [name] RoomJoinResp name
         * @property {string|null} [owner] RoomJoinResp owner
         * @property {number|null} [size] RoomJoinResp size
         * @property {number|null} [max] RoomJoinResp max
         * @property {Array.<string>|null} [members] RoomJoinResp members
         */

        /**
         * Constructs a new RoomJoinResp.
         * @memberof ws
         * @classdesc Represents a RoomJoinResp.
         * @implements IRoomJoinResp
         * @constructor
         * @param {ws.IRoomJoinResp=} [properties] Properties to set
         */
        function RoomJoinResp(properties) {
            this.members = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomJoinResp roomId.
         * @member {string} roomId
         * @memberof ws.RoomJoinResp
         * @instance
         */
        RoomJoinResp.prototype.roomId = "";

        /**
         * RoomJoinResp name.
         * @member {string} name
         * @memberof ws.RoomJoinResp
         * @instance
         */
        RoomJoinResp.prototype.name = "";

        /**
         * RoomJoinResp owner.
         * @member {string} owner
         * @memberof ws.RoomJoinResp
         * @instance
         */
        RoomJoinResp.prototype.owner = "";

        /**
         * RoomJoinResp size.
         * @member {number} size
         * @memberof ws.RoomJoinResp
         * @instance
         */
        RoomJoinResp.prototype.size = 0;

        /**
         * RoomJoinResp max.
         * @member {number} max
         * @memberof ws.RoomJoinResp
         * @instance
         */
        RoomJoinResp.prototype.max = 0;

        /**
         * RoomJoinResp members.
         * @member {Array.<string>} members
         * @memberof ws.RoomJoinResp
         * @instance
         */
        RoomJoinResp.prototype.members = $util.emptyArray;

        /**
         * Creates a new RoomJoinResp instance using the specified properties.
         * @function create
         * @memberof ws.RoomJoinResp
         * @static
         * @param {ws.IRoomJoinResp=} [properties] Properties to set
         * @returns {ws.RoomJoinResp} RoomJoinResp instance
         */
        RoomJoinResp.create = function create(properties) {
            return new RoomJoinResp(properties);
        };

        /**
         * Encodes the specified RoomJoinResp message. Does not implicitly {@link ws.RoomJoinResp.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomJoinResp
         * @static
         * @param {ws.IRoomJoinResp} message RoomJoinResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomJoinResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.name);
            if (message.owner != null && Object.hasOwnProperty.call(message, "owner"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.owner);
            if (message.size != null && Object.hasOwnProperty.call(message, "size"))
                writer.uint32(/* id 4, wireType 0 =*/32).int32(message.size);
            if (message.max != null && Object.hasOwnProperty.call(message, "max"))
                writer.uint32(/* id 5, wireType 0 =*/40).int32(message.max);
            if (message.members != null && message.members.length)
                for (let i = 0; i < message.members.length; ++i)
                    writer.uint32(/* id 6, wireType 2 =*/50).string(message.members[i]);
            return writer;
        };

        /**
         * Encodes the specified RoomJoinResp message, length delimited. Does not implicitly {@link ws.RoomJoinResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomJoinResp
         * @static
         * @param {ws.IRoomJoinResp} message RoomJoinResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomJoinResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomJoinResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomJoinResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomJoinResp} RoomJoinResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomJoinResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomJoinResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                case 2:
                    message.name = reader.string();
                    break;
                case 3:
                    message.owner = reader.string();
                    break;
                case 4:
                    message.size = reader.int32();
                    break;
                case 5:
                    message.max = reader.int32();
                    break;
                case 6:
                    if (!(message.members && message.members.length))
                        message.members = [];
                    message.members.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomJoinResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomJoinResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomJoinResp} RoomJoinResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomJoinResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomJoinResp message.
         * @function verify
         * @memberof ws.RoomJoinResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomJoinResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            if (message.owner != null && message.hasOwnProperty("owner"))
                if (!$util.isString(message.owner))
                    return "owner: string expected";
            if (message.size != null && message.hasOwnProperty("size"))
                if (!$util.isInteger(message.size))
                    return "size: integer expected";
            if (message.max != null && message.hasOwnProperty("max"))
                if (!$util.isInteger(message.max))
                    return "max: integer expected";
            if (message.members != null && message.hasOwnProperty("members")) {
                if (!Array.isArray(message.members))
                    return "members: array expected";
                for (let i = 0; i < message.members.length; ++i)
                    if (!$util.isString(message.members[i]))
                        return "members: string[] expected";
            }
            return null;
        };

        /**
         * Creates a RoomJoinResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomJoinResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomJoinResp} RoomJoinResp
         */
        RoomJoinResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomJoinResp)
                return object;
            let message = new $root.ws.RoomJoinResp();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            if (object.name != null)
                message.name = String(object.name);
            if (object.owner != null)
                message.owner = String(object.owner);
            if (object.size != null)
                message.size = object.size | 0;
            if (object.max != null)
                message.max = object.max | 0;
            if (object.members) {
                if (!Array.isArray(object.members))
                    throw TypeError(".ws.RoomJoinResp.members: array expected");
                message.members = [];
                for (let i = 0; i < object.members.length; ++i)
                    message.members[i] = String(object.members[i]);
            }
            return message;
        };

        /**
         * Creates a plain object from a RoomJoinResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomJoinResp
         * @static
         * @param {ws.RoomJoinResp} message RoomJoinResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomJoinResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.members = [];
            if (options.defaults) {
                object.roomId = "";
                object.name = "";
                object.owner = "";
                object.size = 0;
                object.max = 0;
            }
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            if (message.owner != null && message.hasOwnProperty("owner"))
                object.owner = message.owner;
            if (message.size != null && message.hasOwnProperty("size"))
                object.size = message.size;
            if (message.max != null && message.hasOwnProperty("max"))
                object.max = message.max;
            if (message.members && message.members.length) {
                object.members = [];
                for (let j = 0; j < message.members.length; ++j)
                    object.members[j] = message.members[j];
            }
            return object;
        };

        /**
         * Converts this RoomJoinResp to JSON.
         * @function toJSON
         * @memberof ws.RoomJoinResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomJoinResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomJoinResp;
    })();

    ws.RoomLeave = (function() {

        /**
         * Properties of a RoomLeave.
         * @memberof ws
         * @interface IRoomLeave
         * @property {string|null} [roomId] RoomLeave roomId
         */

        /**
         * Constructs a new RoomLeave.
         * @memberof ws
         * @classdesc Represents a RoomLeave.
         * @implements IRoomLeave
         * @constructor
         * @param {ws.IRoomLeave=} [properties] Properties to set
         */
        function RoomLeave(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomLeave roomId.
         * @member {string} roomId
         * @memberof ws.RoomLeave
         * @instance
         */
        RoomLeave.prototype.roomId = "";

        /**
         * Creates a new RoomLeave instance using the specified properties.
         * @function create
         * @memberof ws.RoomLeave
         * @static
         * @param {ws.IRoomLeave=} [properties] Properties to set
         * @returns {ws.RoomLeave} RoomLeave instance
         */
        RoomLeave.create = function create(properties) {
            return new RoomLeave(properties);
        };

        /**
         * Encodes the specified RoomLeave message. Does not implicitly {@link ws.RoomLeave.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomLeave
         * @static
         * @param {ws.IRoomLeave} message RoomLeave message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomLeave.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            return writer;
        };

        /**
         * Encodes the specified RoomLeave message, length delimited. Does not implicitly {@link ws.RoomLeave.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomLeave
         * @static
         * @param {ws.IRoomLeave} message RoomLeave message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomLeave.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomLeave message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomLeave
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomLeave} RoomLeave
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomLeave.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomLeave();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomLeave message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomLeave
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomLeave} RoomLeave
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomLeave.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomLeave message.
         * @function verify
         * @memberof ws.RoomLeave
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomLeave.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            return null;
        };

        /**
         * Creates a RoomLeave message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomLeave
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomLeave} RoomLeave
         */
        RoomLeave.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomLeave)
                return object;
            let message = new $root.ws.RoomLeave();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            return message;
        };

        /**
         * Creates a plain object from a RoomLeave message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomLeave
         * @static
         * @param {ws.RoomLeave} message RoomLeave
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomLeave.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.roomId = "";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            return object;
        };

        /**
         * Converts this RoomLeave to JSON.
         * @function toJSON
         * @memberof ws.RoomLeave
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomLeave.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomLeave;
    })();

    ws.RoomKick = (function() {

        /**
         * Properties of a RoomKick.
         * @memberof ws
         * @interface IRoomKick
         * @property {string|null} [roomId] RoomKick roomId
         * @property {string|null} [targetUid] RoomKick targetUid
         */

        /**
         * Constructs a new RoomKick.
         * @memberof ws
         * @classdesc Represents a RoomKick.
         * @implements IRoomKick
         * @constructor
         * @param {ws.IRoomKick=} [properties] Properties to set
         */
        function RoomKick(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomKick roomId.
         * @member {string} roomId
         * @memberof ws.RoomKick
         * @instance
         */
        RoomKick.prototype.roomId = "";

        /**
         * RoomKick targetUid.
         * @member {string} targetUid
         * @memberof ws.RoomKick
         * @instance
         */
        RoomKick.prototype.targetUid = "";

        /**
         * Creates a new RoomKick instance using the specified properties.
         * @function create
         * @memberof ws.RoomKick
         * @static
         * @param {ws.IRoomKick=} [properties] Properties to set
         * @returns {ws.RoomKick} RoomKick instance
         */
        RoomKick.create = function create(properties) {
            return new RoomKick(properties);
        };

        /**
         * Encodes the specified RoomKick message. Does not implicitly {@link ws.RoomKick.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomKick
         * @static
         * @param {ws.IRoomKick} message RoomKick message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomKick.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            if (message.targetUid != null && Object.hasOwnProperty.call(message, "targetUid"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.targetUid);
            return writer;
        };

        /**
         * Encodes the specified RoomKick message, length delimited. Does not implicitly {@link ws.RoomKick.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomKick
         * @static
         * @param {ws.IRoomKick} message RoomKick message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomKick.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomKick message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomKick
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomKick} RoomKick
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomKick.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomKick();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                case 2:
                    message.targetUid = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomKick message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomKick
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomKick} RoomKick
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomKick.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomKick message.
         * @function verify
         * @memberof ws.RoomKick
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomKick.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            if (message.targetUid != null && message.hasOwnProperty("targetUid"))
                if (!$util.isString(message.targetUid))
                    return "targetUid: string expected";
            return null;
        };

        /**
         * Creates a RoomKick message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomKick
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomKick} RoomKick
         */
        RoomKick.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomKick)
                return object;
            let message = new $root.ws.RoomKick();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            if (object.targetUid != null)
                message.targetUid = String(object.targetUid);
            return message;
        };

        /**
         * Creates a plain object from a RoomKick message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomKick
         * @static
         * @param {ws.RoomKick} message RoomKick
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomKick.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.roomId = "";
                object.targetUid = "";
            }
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            if (message.targetUid != null && message.hasOwnProperty("targetUid"))
                object.targetUid = message.targetUid;
            return object;
        };

        /**
         * Converts this RoomKick to JSON.
         * @function toJSON
         * @memberof ws.RoomKick
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomKick.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomKick;
    })();

    ws.RoomDelete = (function() {

        /**
         * Properties of a RoomDelete.
         * @memberof ws
         * @interface IRoomDelete
         * @property {string|null} [roomId] RoomDelete roomId
         */

        /**
         * Constructs a new RoomDelete.
         * @memberof ws
         * @classdesc Represents a RoomDelete.
         * @implements IRoomDelete
         * @constructor
         * @param {ws.IRoomDelete=} [properties] Properties to set
         */
        function RoomDelete(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomDelete roomId.
         * @member {string} roomId
         * @memberof ws.RoomDelete
         * @instance
         */
        RoomDelete.prototype.roomId = "";

        /**
         * Creates a new RoomDelete instance using the specified properties.
         * @function create
         * @memberof ws.RoomDelete
         * @static
         * @param {ws.IRoomDelete=} [properties] Properties to set
         * @returns {ws.RoomDelete} RoomDelete instance
         */
        RoomDelete.create = function create(properties) {
            return new RoomDelete(properties);
        };

        /**
         * Encodes the specified RoomDelete message. Does not implicitly {@link ws.RoomDelete.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomDelete
         * @static
         * @param {ws.IRoomDelete} message RoomDelete message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomDelete.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            return writer;
        };

        /**
         * Encodes the specified RoomDelete message, length delimited. Does not implicitly {@link ws.RoomDelete.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomDelete
         * @static
         * @param {ws.IRoomDelete} message RoomDelete message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomDelete.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomDelete message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomDelete
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomDelete} RoomDelete
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomDelete.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomDelete();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomDelete message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomDelete
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomDelete} RoomDelete
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomDelete.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomDelete message.
         * @function verify
         * @memberof ws.RoomDelete
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomDelete.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            return null;
        };

        /**
         * Creates a RoomDelete message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomDelete
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomDelete} RoomDelete
         */
        RoomDelete.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomDelete)
                return object;
            let message = new $root.ws.RoomDelete();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            return message;
        };

        /**
         * Creates a plain object from a RoomDelete message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomDelete
         * @static
         * @param {ws.RoomDelete} message RoomDelete
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomDelete.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults)
                object.roomId = "";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            return object;
        };

        /**
         * Converts this RoomDelete to JSON.
         * @function toJSON
         * @memberof ws.RoomDelete
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomDelete.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomDelete;
    })();

    ws.RoomSend = (function() {

        /**
         * Properties of a RoomSend.
         * @memberof ws
         * @interface IRoomSend
         * @property {string|null} [roomId] RoomSend roomId
         * @property {string|null} [content] RoomSend content
         */

        /**
         * Constructs a new RoomSend.
         * @memberof ws
         * @classdesc Represents a RoomSend.
         * @implements IRoomSend
         * @constructor
         * @param {ws.IRoomSend=} [properties] Properties to set
         */
        function RoomSend(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomSend roomId.
         * @member {string} roomId
         * @memberof ws.RoomSend
         * @instance
         */
        RoomSend.prototype.roomId = "";

        /**
         * RoomSend content.
         * @member {string} content
         * @memberof ws.RoomSend
         * @instance
         */
        RoomSend.prototype.content = "";

        /**
         * Creates a new RoomSend instance using the specified properties.
         * @function create
         * @memberof ws.RoomSend
         * @static
         * @param {ws.IRoomSend=} [properties] Properties to set
         * @returns {ws.RoomSend} RoomSend instance
         */
        RoomSend.create = function create(properties) {
            return new RoomSend(properties);
        };

        /**
         * Encodes the specified RoomSend message. Does not implicitly {@link ws.RoomSend.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomSend
         * @static
         * @param {ws.IRoomSend} message RoomSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomSend.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.roomId);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.content);
            return writer;
        };

        /**
         * Encodes the specified RoomSend message, length delimited. Does not implicitly {@link ws.RoomSend.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomSend
         * @static
         * @param {ws.IRoomSend} message RoomSend message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomSend.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomSend message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomSend} RoomSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomSend.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomSend();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.roomId = reader.string();
                    break;
                case 2:
                    message.content = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomSend message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomSend
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomSend} RoomSend
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomSend.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomSend message.
         * @function verify
         * @memberof ws.RoomSend
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomSend.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            return null;
        };

        /**
         * Creates a RoomSend message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomSend
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomSend} RoomSend
         */
        RoomSend.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomSend)
                return object;
            let message = new $root.ws.RoomSend();
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            if (object.content != null)
                message.content = String(object.content);
            return message;
        };

        /**
         * Creates a plain object from a RoomSend message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomSend
         * @static
         * @param {ws.RoomSend} message RoomSend
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomSend.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.roomId = "";
                object.content = "";
            }
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            return object;
        };

        /**
         * Converts this RoomSend to JSON.
         * @function toJSON
         * @memberof ws.RoomSend
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomSend.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomSend;
    })();

    ws.RoomChatNew = (function() {

        /**
         * Properties of a RoomChatNew.
         * @memberof ws
         * @interface IRoomChatNew
         * @property {string|null} [msgId] RoomChatNew msgId
         * @property {string|null} [roomId] RoomChatNew roomId
         * @property {string|null} [fromUid] RoomChatNew fromUid
         * @property {string|null} [fromName] RoomChatNew fromName
         * @property {string|null} [content] RoomChatNew content
         * @property {number|Long|null} [timestamp] RoomChatNew timestamp
         */

        /**
         * Constructs a new RoomChatNew.
         * @memberof ws
         * @classdesc Represents a RoomChatNew.
         * @implements IRoomChatNew
         * @constructor
         * @param {ws.IRoomChatNew=} [properties] Properties to set
         */
        function RoomChatNew(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomChatNew msgId.
         * @member {string} msgId
         * @memberof ws.RoomChatNew
         * @instance
         */
        RoomChatNew.prototype.msgId = "";

        /**
         * RoomChatNew roomId.
         * @member {string} roomId
         * @memberof ws.RoomChatNew
         * @instance
         */
        RoomChatNew.prototype.roomId = "";

        /**
         * RoomChatNew fromUid.
         * @member {string} fromUid
         * @memberof ws.RoomChatNew
         * @instance
         */
        RoomChatNew.prototype.fromUid = "";

        /**
         * RoomChatNew fromName.
         * @member {string} fromName
         * @memberof ws.RoomChatNew
         * @instance
         */
        RoomChatNew.prototype.fromName = "";

        /**
         * RoomChatNew content.
         * @member {string} content
         * @memberof ws.RoomChatNew
         * @instance
         */
        RoomChatNew.prototype.content = "";

        /**
         * RoomChatNew timestamp.
         * @member {number|Long} timestamp
         * @memberof ws.RoomChatNew
         * @instance
         */
        RoomChatNew.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new RoomChatNew instance using the specified properties.
         * @function create
         * @memberof ws.RoomChatNew
         * @static
         * @param {ws.IRoomChatNew=} [properties] Properties to set
         * @returns {ws.RoomChatNew} RoomChatNew instance
         */
        RoomChatNew.create = function create(properties) {
            return new RoomChatNew(properties);
        };

        /**
         * Encodes the specified RoomChatNew message. Does not implicitly {@link ws.RoomChatNew.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomChatNew
         * @static
         * @param {ws.IRoomChatNew} message RoomChatNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomChatNew.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.msgId != null && Object.hasOwnProperty.call(message, "msgId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.msgId);
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.roomId);
            if (message.fromUid != null && Object.hasOwnProperty.call(message, "fromUid"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.fromUid);
            if (message.fromName != null && Object.hasOwnProperty.call(message, "fromName"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.fromName);
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 5, wireType 2 =*/42).string(message.content);
            if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
                writer.uint32(/* id 6, wireType 0 =*/48).int64(message.timestamp);
            return writer;
        };

        /**
         * Encodes the specified RoomChatNew message, length delimited. Does not implicitly {@link ws.RoomChatNew.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomChatNew
         * @static
         * @param {ws.IRoomChatNew} message RoomChatNew message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomChatNew.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomChatNew message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomChatNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomChatNew} RoomChatNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomChatNew.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomChatNew();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.msgId = reader.string();
                    break;
                case 2:
                    message.roomId = reader.string();
                    break;
                case 3:
                    message.fromUid = reader.string();
                    break;
                case 4:
                    message.fromName = reader.string();
                    break;
                case 5:
                    message.content = reader.string();
                    break;
                case 6:
                    message.timestamp = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomChatNew message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomChatNew
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomChatNew} RoomChatNew
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomChatNew.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomChatNew message.
         * @function verify
         * @memberof ws.RoomChatNew
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomChatNew.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                if (!$util.isString(message.msgId))
                    return "msgId: string expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                if (!$util.isString(message.fromUid))
                    return "fromUid: string expected";
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                if (!$util.isString(message.fromName))
                    return "fromName: string expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                    return "timestamp: integer|Long expected";
            return null;
        };

        /**
         * Creates a RoomChatNew message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomChatNew
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomChatNew} RoomChatNew
         */
        RoomChatNew.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomChatNew)
                return object;
            let message = new $root.ws.RoomChatNew();
            if (object.msgId != null)
                message.msgId = String(object.msgId);
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            if (object.fromUid != null)
                message.fromUid = String(object.fromUid);
            if (object.fromName != null)
                message.fromName = String(object.fromName);
            if (object.content != null)
                message.content = String(object.content);
            if (object.timestamp != null)
                if ($util.Long)
                    (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
                else if (typeof object.timestamp === "string")
                    message.timestamp = parseInt(object.timestamp, 10);
                else if (typeof object.timestamp === "number")
                    message.timestamp = object.timestamp;
                else if (typeof object.timestamp === "object")
                    message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a RoomChatNew message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomChatNew
         * @static
         * @param {ws.RoomChatNew} message RoomChatNew
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomChatNew.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.msgId = "";
                object.roomId = "";
                object.fromUid = "";
                object.fromName = "";
                object.content = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.timestamp = options.longs === String ? "0" : 0;
            }
            if (message.msgId != null && message.hasOwnProperty("msgId"))
                object.msgId = message.msgId;
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            if (message.fromUid != null && message.hasOwnProperty("fromUid"))
                object.fromUid = message.fromUid;
            if (message.fromName != null && message.hasOwnProperty("fromName"))
                object.fromName = message.fromName;
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (typeof message.timestamp === "number")
                    object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
                else
                    object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
            return object;
        };

        /**
         * Converts this RoomChatNew to JSON.
         * @function toJSON
         * @memberof ws.RoomChatNew
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomChatNew.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomChatNew;
    })();

    ws.RoomNotice = (function() {

        /**
         * Properties of a RoomNotice.
         * @memberof ws
         * @interface IRoomNotice
         * @property {string|null} [type] RoomNotice type
         * @property {string|null} [roomId] RoomNotice roomId
         * @property {string|null} [uid] RoomNotice uid
         * @property {string|null} [username] RoomNotice username
         * @property {Array.<string>|null} [members] RoomNotice members
         */

        /**
         * Constructs a new RoomNotice.
         * @memberof ws
         * @classdesc Represents a RoomNotice.
         * @implements IRoomNotice
         * @constructor
         * @param {ws.IRoomNotice=} [properties] Properties to set
         */
        function RoomNotice(properties) {
            this.members = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * RoomNotice type.
         * @member {string} type
         * @memberof ws.RoomNotice
         * @instance
         */
        RoomNotice.prototype.type = "";

        /**
         * RoomNotice roomId.
         * @member {string} roomId
         * @memberof ws.RoomNotice
         * @instance
         */
        RoomNotice.prototype.roomId = "";

        /**
         * RoomNotice uid.
         * @member {string} uid
         * @memberof ws.RoomNotice
         * @instance
         */
        RoomNotice.prototype.uid = "";

        /**
         * RoomNotice username.
         * @member {string} username
         * @memberof ws.RoomNotice
         * @instance
         */
        RoomNotice.prototype.username = "";

        /**
         * RoomNotice members.
         * @member {Array.<string>} members
         * @memberof ws.RoomNotice
         * @instance
         */
        RoomNotice.prototype.members = $util.emptyArray;

        /**
         * Creates a new RoomNotice instance using the specified properties.
         * @function create
         * @memberof ws.RoomNotice
         * @static
         * @param {ws.IRoomNotice=} [properties] Properties to set
         * @returns {ws.RoomNotice} RoomNotice instance
         */
        RoomNotice.create = function create(properties) {
            return new RoomNotice(properties);
        };

        /**
         * Encodes the specified RoomNotice message. Does not implicitly {@link ws.RoomNotice.verify|verify} messages.
         * @function encode
         * @memberof ws.RoomNotice
         * @static
         * @param {ws.IRoomNotice} message RoomNotice message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomNotice.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.type != null && Object.hasOwnProperty.call(message, "type"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.type);
            if (message.roomId != null && Object.hasOwnProperty.call(message, "roomId"))
                writer.uint32(/* id 2, wireType 2 =*/18).string(message.roomId);
            if (message.uid != null && Object.hasOwnProperty.call(message, "uid"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.uid);
            if (message.username != null && Object.hasOwnProperty.call(message, "username"))
                writer.uint32(/* id 4, wireType 2 =*/34).string(message.username);
            if (message.members != null && message.members.length)
                for (let i = 0; i < message.members.length; ++i)
                    writer.uint32(/* id 5, wireType 2 =*/42).string(message.members[i]);
            return writer;
        };

        /**
         * Encodes the specified RoomNotice message, length delimited. Does not implicitly {@link ws.RoomNotice.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.RoomNotice
         * @static
         * @param {ws.IRoomNotice} message RoomNotice message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        RoomNotice.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a RoomNotice message from the specified reader or buffer.
         * @function decode
         * @memberof ws.RoomNotice
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.RoomNotice} RoomNotice
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomNotice.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.RoomNotice();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.type = reader.string();
                    break;
                case 2:
                    message.roomId = reader.string();
                    break;
                case 3:
                    message.uid = reader.string();
                    break;
                case 4:
                    message.username = reader.string();
                    break;
                case 5:
                    if (!(message.members && message.members.length))
                        message.members = [];
                    message.members.push(reader.string());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a RoomNotice message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.RoomNotice
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.RoomNotice} RoomNotice
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        RoomNotice.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a RoomNotice message.
         * @function verify
         * @memberof ws.RoomNotice
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        RoomNotice.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.type != null && message.hasOwnProperty("type"))
                if (!$util.isString(message.type))
                    return "type: string expected";
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                if (!$util.isString(message.roomId))
                    return "roomId: string expected";
            if (message.uid != null && message.hasOwnProperty("uid"))
                if (!$util.isString(message.uid))
                    return "uid: string expected";
            if (message.username != null && message.hasOwnProperty("username"))
                if (!$util.isString(message.username))
                    return "username: string expected";
            if (message.members != null && message.hasOwnProperty("members")) {
                if (!Array.isArray(message.members))
                    return "members: array expected";
                for (let i = 0; i < message.members.length; ++i)
                    if (!$util.isString(message.members[i]))
                        return "members: string[] expected";
            }
            return null;
        };

        /**
         * Creates a RoomNotice message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.RoomNotice
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.RoomNotice} RoomNotice
         */
        RoomNotice.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.RoomNotice)
                return object;
            let message = new $root.ws.RoomNotice();
            if (object.type != null)
                message.type = String(object.type);
            if (object.roomId != null)
                message.roomId = String(object.roomId);
            if (object.uid != null)
                message.uid = String(object.uid);
            if (object.username != null)
                message.username = String(object.username);
            if (object.members) {
                if (!Array.isArray(object.members))
                    throw TypeError(".ws.RoomNotice.members: array expected");
                message.members = [];
                for (let i = 0; i < object.members.length; ++i)
                    message.members[i] = String(object.members[i]);
            }
            return message;
        };

        /**
         * Creates a plain object from a RoomNotice message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.RoomNotice
         * @static
         * @param {ws.RoomNotice} message RoomNotice
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        RoomNotice.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.members = [];
            if (options.defaults) {
                object.type = "";
                object.roomId = "";
                object.uid = "";
                object.username = "";
            }
            if (message.type != null && message.hasOwnProperty("type"))
                object.type = message.type;
            if (message.roomId != null && message.hasOwnProperty("roomId"))
                object.roomId = message.roomId;
            if (message.uid != null && message.hasOwnProperty("uid"))
                object.uid = message.uid;
            if (message.username != null && message.hasOwnProperty("username"))
                object.username = message.username;
            if (message.members && message.members.length) {
                object.members = [];
                for (let j = 0; j < message.members.length; ++j)
                    object.members[j] = message.members[j];
            }
            return object;
        };

        /**
         * Converts this RoomNotice to JSON.
         * @function toJSON
         * @memberof ws.RoomNotice
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        RoomNotice.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return RoomNotice;
    })();

    ws.SystemBroadcast = (function() {

        /**
         * Properties of a SystemBroadcast.
         * @memberof ws
         * @interface ISystemBroadcast
         * @property {string|null} [content] SystemBroadcast content
         * @property {number|Long|null} [timestamp] SystemBroadcast timestamp
         */

        /**
         * Constructs a new SystemBroadcast.
         * @memberof ws
         * @classdesc Represents a SystemBroadcast.
         * @implements ISystemBroadcast
         * @constructor
         * @param {ws.ISystemBroadcast=} [properties] Properties to set
         */
        function SystemBroadcast(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * SystemBroadcast content.
         * @member {string} content
         * @memberof ws.SystemBroadcast
         * @instance
         */
        SystemBroadcast.prototype.content = "";

        /**
         * SystemBroadcast timestamp.
         * @member {number|Long} timestamp
         * @memberof ws.SystemBroadcast
         * @instance
         */
        SystemBroadcast.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new SystemBroadcast instance using the specified properties.
         * @function create
         * @memberof ws.SystemBroadcast
         * @static
         * @param {ws.ISystemBroadcast=} [properties] Properties to set
         * @returns {ws.SystemBroadcast} SystemBroadcast instance
         */
        SystemBroadcast.create = function create(properties) {
            return new SystemBroadcast(properties);
        };

        /**
         * Encodes the specified SystemBroadcast message. Does not implicitly {@link ws.SystemBroadcast.verify|verify} messages.
         * @function encode
         * @memberof ws.SystemBroadcast
         * @static
         * @param {ws.ISystemBroadcast} message SystemBroadcast message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SystemBroadcast.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.content);
            if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.timestamp);
            return writer;
        };

        /**
         * Encodes the specified SystemBroadcast message, length delimited. Does not implicitly {@link ws.SystemBroadcast.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.SystemBroadcast
         * @static
         * @param {ws.ISystemBroadcast} message SystemBroadcast message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SystemBroadcast.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a SystemBroadcast message from the specified reader or buffer.
         * @function decode
         * @memberof ws.SystemBroadcast
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.SystemBroadcast} SystemBroadcast
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SystemBroadcast.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.SystemBroadcast();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.content = reader.string();
                    break;
                case 2:
                    message.timestamp = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a SystemBroadcast message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.SystemBroadcast
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.SystemBroadcast} SystemBroadcast
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SystemBroadcast.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a SystemBroadcast message.
         * @function verify
         * @memberof ws.SystemBroadcast
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        SystemBroadcast.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                    return "timestamp: integer|Long expected";
            return null;
        };

        /**
         * Creates a SystemBroadcast message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.SystemBroadcast
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.SystemBroadcast} SystemBroadcast
         */
        SystemBroadcast.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.SystemBroadcast)
                return object;
            let message = new $root.ws.SystemBroadcast();
            if (object.content != null)
                message.content = String(object.content);
            if (object.timestamp != null)
                if ($util.Long)
                    (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
                else if (typeof object.timestamp === "string")
                    message.timestamp = parseInt(object.timestamp, 10);
                else if (typeof object.timestamp === "number")
                    message.timestamp = object.timestamp;
                else if (typeof object.timestamp === "object")
                    message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a SystemBroadcast message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.SystemBroadcast
         * @static
         * @param {ws.SystemBroadcast} message SystemBroadcast
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        SystemBroadcast.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.content = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.timestamp = options.longs === String ? "0" : 0;
            }
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (typeof message.timestamp === "number")
                    object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
                else
                    object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
            return object;
        };

        /**
         * Converts this SystemBroadcast to JSON.
         * @function toJSON
         * @memberof ws.SystemBroadcast
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        SystemBroadcast.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return SystemBroadcast;
    })();

    ws.SystemPush = (function() {

        /**
         * Properties of a SystemPush.
         * @memberof ws
         * @interface ISystemPush
         * @property {string|null} [content] SystemPush content
         * @property {number|Long|null} [timestamp] SystemPush timestamp
         */

        /**
         * Constructs a new SystemPush.
         * @memberof ws
         * @classdesc Represents a SystemPush.
         * @implements ISystemPush
         * @constructor
         * @param {ws.ISystemPush=} [properties] Properties to set
         */
        function SystemPush(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * SystemPush content.
         * @member {string} content
         * @memberof ws.SystemPush
         * @instance
         */
        SystemPush.prototype.content = "";

        /**
         * SystemPush timestamp.
         * @member {number|Long} timestamp
         * @memberof ws.SystemPush
         * @instance
         */
        SystemPush.prototype.timestamp = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new SystemPush instance using the specified properties.
         * @function create
         * @memberof ws.SystemPush
         * @static
         * @param {ws.ISystemPush=} [properties] Properties to set
         * @returns {ws.SystemPush} SystemPush instance
         */
        SystemPush.create = function create(properties) {
            return new SystemPush(properties);
        };

        /**
         * Encodes the specified SystemPush message. Does not implicitly {@link ws.SystemPush.verify|verify} messages.
         * @function encode
         * @memberof ws.SystemPush
         * @static
         * @param {ws.ISystemPush} message SystemPush message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SystemPush.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.content != null && Object.hasOwnProperty.call(message, "content"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.content);
            if (message.timestamp != null && Object.hasOwnProperty.call(message, "timestamp"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.timestamp);
            return writer;
        };

        /**
         * Encodes the specified SystemPush message, length delimited. Does not implicitly {@link ws.SystemPush.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.SystemPush
         * @static
         * @param {ws.ISystemPush} message SystemPush message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        SystemPush.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a SystemPush message from the specified reader or buffer.
         * @function decode
         * @memberof ws.SystemPush
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.SystemPush} SystemPush
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SystemPush.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.SystemPush();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.content = reader.string();
                    break;
                case 2:
                    message.timestamp = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a SystemPush message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.SystemPush
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.SystemPush} SystemPush
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        SystemPush.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a SystemPush message.
         * @function verify
         * @memberof ws.SystemPush
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        SystemPush.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.content != null && message.hasOwnProperty("content"))
                if (!$util.isString(message.content))
                    return "content: string expected";
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (!$util.isInteger(message.timestamp) && !(message.timestamp && $util.isInteger(message.timestamp.low) && $util.isInteger(message.timestamp.high)))
                    return "timestamp: integer|Long expected";
            return null;
        };

        /**
         * Creates a SystemPush message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.SystemPush
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.SystemPush} SystemPush
         */
        SystemPush.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.SystemPush)
                return object;
            let message = new $root.ws.SystemPush();
            if (object.content != null)
                message.content = String(object.content);
            if (object.timestamp != null)
                if ($util.Long)
                    (message.timestamp = $util.Long.fromValue(object.timestamp)).unsigned = false;
                else if (typeof object.timestamp === "string")
                    message.timestamp = parseInt(object.timestamp, 10);
                else if (typeof object.timestamp === "number")
                    message.timestamp = object.timestamp;
                else if (typeof object.timestamp === "object")
                    message.timestamp = new $util.LongBits(object.timestamp.low >>> 0, object.timestamp.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from a SystemPush message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.SystemPush
         * @static
         * @param {ws.SystemPush} message SystemPush
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        SystemPush.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.content = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.timestamp = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.timestamp = options.longs === String ? "0" : 0;
            }
            if (message.content != null && message.hasOwnProperty("content"))
                object.content = message.content;
            if (message.timestamp != null && message.hasOwnProperty("timestamp"))
                if (typeof message.timestamp === "number")
                    object.timestamp = options.longs === String ? String(message.timestamp) : message.timestamp;
                else
                    object.timestamp = options.longs === String ? $util.Long.prototype.toString.call(message.timestamp) : options.longs === Number ? new $util.LongBits(message.timestamp.low >>> 0, message.timestamp.high >>> 0).toNumber() : message.timestamp;
            return object;
        };

        /**
         * Converts this SystemPush to JSON.
         * @function toJSON
         * @memberof ws.SystemPush
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        SystemPush.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return SystemPush;
    })();

    ws.InventoryAdd = (function() {

        /**
         * Properties of an InventoryAdd.
         * @memberof ws
         * @interface IInventoryAdd
         * @property {string|null} [itemId] InventoryAdd itemId
         * @property {number|Long|null} [count] InventoryAdd count
         */

        /**
         * Constructs a new InventoryAdd.
         * @memberof ws
         * @classdesc Represents an InventoryAdd.
         * @implements IInventoryAdd
         * @constructor
         * @param {ws.IInventoryAdd=} [properties] Properties to set
         */
        function InventoryAdd(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryAdd itemId.
         * @member {string} itemId
         * @memberof ws.InventoryAdd
         * @instance
         */
        InventoryAdd.prototype.itemId = "";

        /**
         * InventoryAdd count.
         * @member {number|Long} count
         * @memberof ws.InventoryAdd
         * @instance
         */
        InventoryAdd.prototype.count = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new InventoryAdd instance using the specified properties.
         * @function create
         * @memberof ws.InventoryAdd
         * @static
         * @param {ws.IInventoryAdd=} [properties] Properties to set
         * @returns {ws.InventoryAdd} InventoryAdd instance
         */
        InventoryAdd.create = function create(properties) {
            return new InventoryAdd(properties);
        };

        /**
         * Encodes the specified InventoryAdd message. Does not implicitly {@link ws.InventoryAdd.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryAdd
         * @static
         * @param {ws.IInventoryAdd} message InventoryAdd message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryAdd.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.count != null && Object.hasOwnProperty.call(message, "count"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.count);
            return writer;
        };

        /**
         * Encodes the specified InventoryAdd message, length delimited. Does not implicitly {@link ws.InventoryAdd.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryAdd
         * @static
         * @param {ws.IInventoryAdd} message InventoryAdd message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryAdd.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryAdd message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryAdd
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryAdd} InventoryAdd
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryAdd.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryAdd();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.count = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryAdd message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryAdd
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryAdd} InventoryAdd
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryAdd.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryAdd message.
         * @function verify
         * @memberof ws.InventoryAdd
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryAdd.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.count != null && message.hasOwnProperty("count"))
                if (!$util.isInteger(message.count) && !(message.count && $util.isInteger(message.count.low) && $util.isInteger(message.count.high)))
                    return "count: integer|Long expected";
            return null;
        };

        /**
         * Creates an InventoryAdd message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryAdd
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryAdd} InventoryAdd
         */
        InventoryAdd.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryAdd)
                return object;
            let message = new $root.ws.InventoryAdd();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.count != null)
                if ($util.Long)
                    (message.count = $util.Long.fromValue(object.count)).unsigned = false;
                else if (typeof object.count === "string")
                    message.count = parseInt(object.count, 10);
                else if (typeof object.count === "number")
                    message.count = object.count;
                else if (typeof object.count === "object")
                    message.count = new $util.LongBits(object.count.low >>> 0, object.count.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an InventoryAdd message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryAdd
         * @static
         * @param {ws.InventoryAdd} message InventoryAdd
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryAdd.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.count = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.count = options.longs === String ? "0" : 0;
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.count != null && message.hasOwnProperty("count"))
                if (typeof message.count === "number")
                    object.count = options.longs === String ? String(message.count) : message.count;
                else
                    object.count = options.longs === String ? $util.Long.prototype.toString.call(message.count) : options.longs === Number ? new $util.LongBits(message.count.low >>> 0, message.count.high >>> 0).toNumber() : message.count;
            return object;
        };

        /**
         * Converts this InventoryAdd to JSON.
         * @function toJSON
         * @memberof ws.InventoryAdd
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryAdd.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryAdd;
    })();

    ws.InventoryConsume = (function() {

        /**
         * Properties of an InventoryConsume.
         * @memberof ws
         * @interface IInventoryConsume
         * @property {string|null} [itemId] InventoryConsume itemId
         * @property {number|Long|null} [count] InventoryConsume count
         */

        /**
         * Constructs a new InventoryConsume.
         * @memberof ws
         * @classdesc Represents an InventoryConsume.
         * @implements IInventoryConsume
         * @constructor
         * @param {ws.IInventoryConsume=} [properties] Properties to set
         */
        function InventoryConsume(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryConsume itemId.
         * @member {string} itemId
         * @memberof ws.InventoryConsume
         * @instance
         */
        InventoryConsume.prototype.itemId = "";

        /**
         * InventoryConsume count.
         * @member {number|Long} count
         * @memberof ws.InventoryConsume
         * @instance
         */
        InventoryConsume.prototype.count = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new InventoryConsume instance using the specified properties.
         * @function create
         * @memberof ws.InventoryConsume
         * @static
         * @param {ws.IInventoryConsume=} [properties] Properties to set
         * @returns {ws.InventoryConsume} InventoryConsume instance
         */
        InventoryConsume.create = function create(properties) {
            return new InventoryConsume(properties);
        };

        /**
         * Encodes the specified InventoryConsume message. Does not implicitly {@link ws.InventoryConsume.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryConsume
         * @static
         * @param {ws.IInventoryConsume} message InventoryConsume message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryConsume.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.count != null && Object.hasOwnProperty.call(message, "count"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.count);
            return writer;
        };

        /**
         * Encodes the specified InventoryConsume message, length delimited. Does not implicitly {@link ws.InventoryConsume.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryConsume
         * @static
         * @param {ws.IInventoryConsume} message InventoryConsume message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryConsume.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryConsume message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryConsume
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryConsume} InventoryConsume
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryConsume.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryConsume();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.count = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryConsume message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryConsume
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryConsume} InventoryConsume
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryConsume.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryConsume message.
         * @function verify
         * @memberof ws.InventoryConsume
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryConsume.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.count != null && message.hasOwnProperty("count"))
                if (!$util.isInteger(message.count) && !(message.count && $util.isInteger(message.count.low) && $util.isInteger(message.count.high)))
                    return "count: integer|Long expected";
            return null;
        };

        /**
         * Creates an InventoryConsume message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryConsume
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryConsume} InventoryConsume
         */
        InventoryConsume.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryConsume)
                return object;
            let message = new $root.ws.InventoryConsume();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.count != null)
                if ($util.Long)
                    (message.count = $util.Long.fromValue(object.count)).unsigned = false;
                else if (typeof object.count === "string")
                    message.count = parseInt(object.count, 10);
                else if (typeof object.count === "number")
                    message.count = object.count;
                else if (typeof object.count === "object")
                    message.count = new $util.LongBits(object.count.low >>> 0, object.count.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an InventoryConsume message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryConsume
         * @static
         * @param {ws.InventoryConsume} message InventoryConsume
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryConsume.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.count = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.count = options.longs === String ? "0" : 0;
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.count != null && message.hasOwnProperty("count"))
                if (typeof message.count === "number")
                    object.count = options.longs === String ? String(message.count) : message.count;
                else
                    object.count = options.longs === String ? $util.Long.prototype.toString.call(message.count) : options.longs === Number ? new $util.LongBits(message.count.low >>> 0, message.count.high >>> 0).toNumber() : message.count;
            return object;
        };

        /**
         * Converts this InventoryConsume to JSON.
         * @function toJSON
         * @memberof ws.InventoryConsume
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryConsume.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryConsume;
    })();

    ws.InventoryDrop = (function() {

        /**
         * Properties of an InventoryDrop.
         * @memberof ws
         * @interface IInventoryDrop
         * @property {string|null} [itemId] InventoryDrop itemId
         * @property {number|Long|null} [count] InventoryDrop count
         */

        /**
         * Constructs a new InventoryDrop.
         * @memberof ws
         * @classdesc Represents an InventoryDrop.
         * @implements IInventoryDrop
         * @constructor
         * @param {ws.IInventoryDrop=} [properties] Properties to set
         */
        function InventoryDrop(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryDrop itemId.
         * @member {string} itemId
         * @memberof ws.InventoryDrop
         * @instance
         */
        InventoryDrop.prototype.itemId = "";

        /**
         * InventoryDrop count.
         * @member {number|Long} count
         * @memberof ws.InventoryDrop
         * @instance
         */
        InventoryDrop.prototype.count = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new InventoryDrop instance using the specified properties.
         * @function create
         * @memberof ws.InventoryDrop
         * @static
         * @param {ws.IInventoryDrop=} [properties] Properties to set
         * @returns {ws.InventoryDrop} InventoryDrop instance
         */
        InventoryDrop.create = function create(properties) {
            return new InventoryDrop(properties);
        };

        /**
         * Encodes the specified InventoryDrop message. Does not implicitly {@link ws.InventoryDrop.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryDrop
         * @static
         * @param {ws.IInventoryDrop} message InventoryDrop message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryDrop.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.count != null && Object.hasOwnProperty.call(message, "count"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.count);
            return writer;
        };

        /**
         * Encodes the specified InventoryDrop message, length delimited. Does not implicitly {@link ws.InventoryDrop.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryDrop
         * @static
         * @param {ws.IInventoryDrop} message InventoryDrop message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryDrop.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryDrop message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryDrop
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryDrop} InventoryDrop
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryDrop.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryDrop();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.count = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryDrop message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryDrop
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryDrop} InventoryDrop
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryDrop.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryDrop message.
         * @function verify
         * @memberof ws.InventoryDrop
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryDrop.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.count != null && message.hasOwnProperty("count"))
                if (!$util.isInteger(message.count) && !(message.count && $util.isInteger(message.count.low) && $util.isInteger(message.count.high)))
                    return "count: integer|Long expected";
            return null;
        };

        /**
         * Creates an InventoryDrop message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryDrop
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryDrop} InventoryDrop
         */
        InventoryDrop.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryDrop)
                return object;
            let message = new $root.ws.InventoryDrop();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.count != null)
                if ($util.Long)
                    (message.count = $util.Long.fromValue(object.count)).unsigned = false;
                else if (typeof object.count === "string")
                    message.count = parseInt(object.count, 10);
                else if (typeof object.count === "number")
                    message.count = object.count;
                else if (typeof object.count === "object")
                    message.count = new $util.LongBits(object.count.low >>> 0, object.count.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an InventoryDrop message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryDrop
         * @static
         * @param {ws.InventoryDrop} message InventoryDrop
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryDrop.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.count = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.count = options.longs === String ? "0" : 0;
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.count != null && message.hasOwnProperty("count"))
                if (typeof message.count === "number")
                    object.count = options.longs === String ? String(message.count) : message.count;
                else
                    object.count = options.longs === String ? $util.Long.prototype.toString.call(message.count) : options.longs === Number ? new $util.LongBits(message.count.low >>> 0, message.count.high >>> 0).toNumber() : message.count;
            return object;
        };

        /**
         * Converts this InventoryDrop to JSON.
         * @function toJSON
         * @memberof ws.InventoryDrop
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryDrop.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryDrop;
    })();

    ws.InventoryGet = (function() {

        /**
         * Properties of an InventoryGet.
         * @memberof ws
         * @interface IInventoryGet
         */

        /**
         * Constructs a new InventoryGet.
         * @memberof ws
         * @classdesc Represents an InventoryGet.
         * @implements IInventoryGet
         * @constructor
         * @param {ws.IInventoryGet=} [properties] Properties to set
         */
        function InventoryGet(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * Creates a new InventoryGet instance using the specified properties.
         * @function create
         * @memberof ws.InventoryGet
         * @static
         * @param {ws.IInventoryGet=} [properties] Properties to set
         * @returns {ws.InventoryGet} InventoryGet instance
         */
        InventoryGet.create = function create(properties) {
            return new InventoryGet(properties);
        };

        /**
         * Encodes the specified InventoryGet message. Does not implicitly {@link ws.InventoryGet.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryGet
         * @static
         * @param {ws.IInventoryGet} message InventoryGet message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryGet.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            return writer;
        };

        /**
         * Encodes the specified InventoryGet message, length delimited. Does not implicitly {@link ws.InventoryGet.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryGet
         * @static
         * @param {ws.IInventoryGet} message InventoryGet message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryGet.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryGet message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryGet
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryGet} InventoryGet
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryGet.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryGet();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryGet message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryGet
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryGet} InventoryGet
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryGet.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryGet message.
         * @function verify
         * @memberof ws.InventoryGet
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryGet.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            return null;
        };

        /**
         * Creates an InventoryGet message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryGet
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryGet} InventoryGet
         */
        InventoryGet.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryGet)
                return object;
            return new $root.ws.InventoryGet();
        };

        /**
         * Creates a plain object from an InventoryGet message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryGet
         * @static
         * @param {ws.InventoryGet} message InventoryGet
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryGet.toObject = function toObject() {
            return {};
        };

        /**
         * Converts this InventoryGet to JSON.
         * @function toJSON
         * @memberof ws.InventoryGet
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryGet.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryGet;
    })();

    ws.InventoryAddResp = (function() {

        /**
         * Properties of an InventoryAddResp.
         * @memberof ws
         * @interface IInventoryAddResp
         * @property {string|null} [itemId] InventoryAddResp itemId
         * @property {number|Long|null} [newCount] InventoryAddResp newCount
         */

        /**
         * Constructs a new InventoryAddResp.
         * @memberof ws
         * @classdesc Represents an InventoryAddResp.
         * @implements IInventoryAddResp
         * @constructor
         * @param {ws.IInventoryAddResp=} [properties] Properties to set
         */
        function InventoryAddResp(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryAddResp itemId.
         * @member {string} itemId
         * @memberof ws.InventoryAddResp
         * @instance
         */
        InventoryAddResp.prototype.itemId = "";

        /**
         * InventoryAddResp newCount.
         * @member {number|Long} newCount
         * @memberof ws.InventoryAddResp
         * @instance
         */
        InventoryAddResp.prototype.newCount = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new InventoryAddResp instance using the specified properties.
         * @function create
         * @memberof ws.InventoryAddResp
         * @static
         * @param {ws.IInventoryAddResp=} [properties] Properties to set
         * @returns {ws.InventoryAddResp} InventoryAddResp instance
         */
        InventoryAddResp.create = function create(properties) {
            return new InventoryAddResp(properties);
        };

        /**
         * Encodes the specified InventoryAddResp message. Does not implicitly {@link ws.InventoryAddResp.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryAddResp
         * @static
         * @param {ws.IInventoryAddResp} message InventoryAddResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryAddResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.newCount != null && Object.hasOwnProperty.call(message, "newCount"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.newCount);
            return writer;
        };

        /**
         * Encodes the specified InventoryAddResp message, length delimited. Does not implicitly {@link ws.InventoryAddResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryAddResp
         * @static
         * @param {ws.IInventoryAddResp} message InventoryAddResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryAddResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryAddResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryAddResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryAddResp} InventoryAddResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryAddResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryAddResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.newCount = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryAddResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryAddResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryAddResp} InventoryAddResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryAddResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryAddResp message.
         * @function verify
         * @memberof ws.InventoryAddResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryAddResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.newCount != null && message.hasOwnProperty("newCount"))
                if (!$util.isInteger(message.newCount) && !(message.newCount && $util.isInteger(message.newCount.low) && $util.isInteger(message.newCount.high)))
                    return "newCount: integer|Long expected";
            return null;
        };

        /**
         * Creates an InventoryAddResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryAddResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryAddResp} InventoryAddResp
         */
        InventoryAddResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryAddResp)
                return object;
            let message = new $root.ws.InventoryAddResp();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.newCount != null)
                if ($util.Long)
                    (message.newCount = $util.Long.fromValue(object.newCount)).unsigned = false;
                else if (typeof object.newCount === "string")
                    message.newCount = parseInt(object.newCount, 10);
                else if (typeof object.newCount === "number")
                    message.newCount = object.newCount;
                else if (typeof object.newCount === "object")
                    message.newCount = new $util.LongBits(object.newCount.low >>> 0, object.newCount.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an InventoryAddResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryAddResp
         * @static
         * @param {ws.InventoryAddResp} message InventoryAddResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryAddResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.newCount = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.newCount = options.longs === String ? "0" : 0;
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.newCount != null && message.hasOwnProperty("newCount"))
                if (typeof message.newCount === "number")
                    object.newCount = options.longs === String ? String(message.newCount) : message.newCount;
                else
                    object.newCount = options.longs === String ? $util.Long.prototype.toString.call(message.newCount) : options.longs === Number ? new $util.LongBits(message.newCount.low >>> 0, message.newCount.high >>> 0).toNumber() : message.newCount;
            return object;
        };

        /**
         * Converts this InventoryAddResp to JSON.
         * @function toJSON
         * @memberof ws.InventoryAddResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryAddResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryAddResp;
    })();

    ws.InventoryConsumeResp = (function() {

        /**
         * Properties of an InventoryConsumeResp.
         * @memberof ws
         * @interface IInventoryConsumeResp
         * @property {string|null} [itemId] InventoryConsumeResp itemId
         * @property {number|Long|null} [newCount] InventoryConsumeResp newCount
         */

        /**
         * Constructs a new InventoryConsumeResp.
         * @memberof ws
         * @classdesc Represents an InventoryConsumeResp.
         * @implements IInventoryConsumeResp
         * @constructor
         * @param {ws.IInventoryConsumeResp=} [properties] Properties to set
         */
        function InventoryConsumeResp(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryConsumeResp itemId.
         * @member {string} itemId
         * @memberof ws.InventoryConsumeResp
         * @instance
         */
        InventoryConsumeResp.prototype.itemId = "";

        /**
         * InventoryConsumeResp newCount.
         * @member {number|Long} newCount
         * @memberof ws.InventoryConsumeResp
         * @instance
         */
        InventoryConsumeResp.prototype.newCount = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new InventoryConsumeResp instance using the specified properties.
         * @function create
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {ws.IInventoryConsumeResp=} [properties] Properties to set
         * @returns {ws.InventoryConsumeResp} InventoryConsumeResp instance
         */
        InventoryConsumeResp.create = function create(properties) {
            return new InventoryConsumeResp(properties);
        };

        /**
         * Encodes the specified InventoryConsumeResp message. Does not implicitly {@link ws.InventoryConsumeResp.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {ws.IInventoryConsumeResp} message InventoryConsumeResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryConsumeResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.newCount != null && Object.hasOwnProperty.call(message, "newCount"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.newCount);
            return writer;
        };

        /**
         * Encodes the specified InventoryConsumeResp message, length delimited. Does not implicitly {@link ws.InventoryConsumeResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {ws.IInventoryConsumeResp} message InventoryConsumeResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryConsumeResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryConsumeResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryConsumeResp} InventoryConsumeResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryConsumeResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryConsumeResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.newCount = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryConsumeResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryConsumeResp} InventoryConsumeResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryConsumeResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryConsumeResp message.
         * @function verify
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryConsumeResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.newCount != null && message.hasOwnProperty("newCount"))
                if (!$util.isInteger(message.newCount) && !(message.newCount && $util.isInteger(message.newCount.low) && $util.isInteger(message.newCount.high)))
                    return "newCount: integer|Long expected";
            return null;
        };

        /**
         * Creates an InventoryConsumeResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryConsumeResp} InventoryConsumeResp
         */
        InventoryConsumeResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryConsumeResp)
                return object;
            let message = new $root.ws.InventoryConsumeResp();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.newCount != null)
                if ($util.Long)
                    (message.newCount = $util.Long.fromValue(object.newCount)).unsigned = false;
                else if (typeof object.newCount === "string")
                    message.newCount = parseInt(object.newCount, 10);
                else if (typeof object.newCount === "number")
                    message.newCount = object.newCount;
                else if (typeof object.newCount === "object")
                    message.newCount = new $util.LongBits(object.newCount.low >>> 0, object.newCount.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an InventoryConsumeResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryConsumeResp
         * @static
         * @param {ws.InventoryConsumeResp} message InventoryConsumeResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryConsumeResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.newCount = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.newCount = options.longs === String ? "0" : 0;
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.newCount != null && message.hasOwnProperty("newCount"))
                if (typeof message.newCount === "number")
                    object.newCount = options.longs === String ? String(message.newCount) : message.newCount;
                else
                    object.newCount = options.longs === String ? $util.Long.prototype.toString.call(message.newCount) : options.longs === Number ? new $util.LongBits(message.newCount.low >>> 0, message.newCount.high >>> 0).toNumber() : message.newCount;
            return object;
        };

        /**
         * Converts this InventoryConsumeResp to JSON.
         * @function toJSON
         * @memberof ws.InventoryConsumeResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryConsumeResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryConsumeResp;
    })();

    ws.InventoryDropResp = (function() {

        /**
         * Properties of an InventoryDropResp.
         * @memberof ws
         * @interface IInventoryDropResp
         * @property {string|null} [itemId] InventoryDropResp itemId
         * @property {number|Long|null} [newCount] InventoryDropResp newCount
         */

        /**
         * Constructs a new InventoryDropResp.
         * @memberof ws
         * @classdesc Represents an InventoryDropResp.
         * @implements IInventoryDropResp
         * @constructor
         * @param {ws.IInventoryDropResp=} [properties] Properties to set
         */
        function InventoryDropResp(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryDropResp itemId.
         * @member {string} itemId
         * @memberof ws.InventoryDropResp
         * @instance
         */
        InventoryDropResp.prototype.itemId = "";

        /**
         * InventoryDropResp newCount.
         * @member {number|Long} newCount
         * @memberof ws.InventoryDropResp
         * @instance
         */
        InventoryDropResp.prototype.newCount = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * Creates a new InventoryDropResp instance using the specified properties.
         * @function create
         * @memberof ws.InventoryDropResp
         * @static
         * @param {ws.IInventoryDropResp=} [properties] Properties to set
         * @returns {ws.InventoryDropResp} InventoryDropResp instance
         */
        InventoryDropResp.create = function create(properties) {
            return new InventoryDropResp(properties);
        };

        /**
         * Encodes the specified InventoryDropResp message. Does not implicitly {@link ws.InventoryDropResp.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryDropResp
         * @static
         * @param {ws.IInventoryDropResp} message InventoryDropResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryDropResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.newCount != null && Object.hasOwnProperty.call(message, "newCount"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.newCount);
            return writer;
        };

        /**
         * Encodes the specified InventoryDropResp message, length delimited. Does not implicitly {@link ws.InventoryDropResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryDropResp
         * @static
         * @param {ws.IInventoryDropResp} message InventoryDropResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryDropResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryDropResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryDropResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryDropResp} InventoryDropResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryDropResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryDropResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.newCount = reader.int64();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryDropResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryDropResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryDropResp} InventoryDropResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryDropResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryDropResp message.
         * @function verify
         * @memberof ws.InventoryDropResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryDropResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.newCount != null && message.hasOwnProperty("newCount"))
                if (!$util.isInteger(message.newCount) && !(message.newCount && $util.isInteger(message.newCount.low) && $util.isInteger(message.newCount.high)))
                    return "newCount: integer|Long expected";
            return null;
        };

        /**
         * Creates an InventoryDropResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryDropResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryDropResp} InventoryDropResp
         */
        InventoryDropResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryDropResp)
                return object;
            let message = new $root.ws.InventoryDropResp();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.newCount != null)
                if ($util.Long)
                    (message.newCount = $util.Long.fromValue(object.newCount)).unsigned = false;
                else if (typeof object.newCount === "string")
                    message.newCount = parseInt(object.newCount, 10);
                else if (typeof object.newCount === "number")
                    message.newCount = object.newCount;
                else if (typeof object.newCount === "object")
                    message.newCount = new $util.LongBits(object.newCount.low >>> 0, object.newCount.high >>> 0).toNumber();
            return message;
        };

        /**
         * Creates a plain object from an InventoryDropResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryDropResp
         * @static
         * @param {ws.InventoryDropResp} message InventoryDropResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryDropResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.newCount = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.newCount = options.longs === String ? "0" : 0;
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.newCount != null && message.hasOwnProperty("newCount"))
                if (typeof message.newCount === "number")
                    object.newCount = options.longs === String ? String(message.newCount) : message.newCount;
                else
                    object.newCount = options.longs === String ? $util.Long.prototype.toString.call(message.newCount) : options.longs === Number ? new $util.LongBits(message.newCount.low >>> 0, message.newCount.high >>> 0).toNumber() : message.newCount;
            return object;
        };

        /**
         * Converts this InventoryDropResp to JSON.
         * @function toJSON
         * @memberof ws.InventoryDropResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryDropResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryDropResp;
    })();

    ws.InventoryGetResp = (function() {

        /**
         * Properties of an InventoryGetResp.
         * @memberof ws
         * @interface IInventoryGetResp
         * @property {Array.<ws.IItemEntry>|null} [items] InventoryGetResp items
         * @property {number|null} [total] InventoryGetResp total
         */

        /**
         * Constructs a new InventoryGetResp.
         * @memberof ws
         * @classdesc Represents an InventoryGetResp.
         * @implements IInventoryGetResp
         * @constructor
         * @param {ws.IInventoryGetResp=} [properties] Properties to set
         */
        function InventoryGetResp(properties) {
            this.items = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * InventoryGetResp items.
         * @member {Array.<ws.IItemEntry>} items
         * @memberof ws.InventoryGetResp
         * @instance
         */
        InventoryGetResp.prototype.items = $util.emptyArray;

        /**
         * InventoryGetResp total.
         * @member {number} total
         * @memberof ws.InventoryGetResp
         * @instance
         */
        InventoryGetResp.prototype.total = 0;

        /**
         * Creates a new InventoryGetResp instance using the specified properties.
         * @function create
         * @memberof ws.InventoryGetResp
         * @static
         * @param {ws.IInventoryGetResp=} [properties] Properties to set
         * @returns {ws.InventoryGetResp} InventoryGetResp instance
         */
        InventoryGetResp.create = function create(properties) {
            return new InventoryGetResp(properties);
        };

        /**
         * Encodes the specified InventoryGetResp message. Does not implicitly {@link ws.InventoryGetResp.verify|verify} messages.
         * @function encode
         * @memberof ws.InventoryGetResp
         * @static
         * @param {ws.IInventoryGetResp} message InventoryGetResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryGetResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.items != null && message.items.length)
                for (let i = 0; i < message.items.length; ++i)
                    $root.ws.ItemEntry.encode(message.items[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            if (message.total != null && Object.hasOwnProperty.call(message, "total"))
                writer.uint32(/* id 2, wireType 0 =*/16).int32(message.total);
            return writer;
        };

        /**
         * Encodes the specified InventoryGetResp message, length delimited. Does not implicitly {@link ws.InventoryGetResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.InventoryGetResp
         * @static
         * @param {ws.IInventoryGetResp} message InventoryGetResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        InventoryGetResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an InventoryGetResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.InventoryGetResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.InventoryGetResp} InventoryGetResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryGetResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.InventoryGetResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    if (!(message.items && message.items.length))
                        message.items = [];
                    message.items.push($root.ws.ItemEntry.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.total = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an InventoryGetResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.InventoryGetResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.InventoryGetResp} InventoryGetResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        InventoryGetResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an InventoryGetResp message.
         * @function verify
         * @memberof ws.InventoryGetResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        InventoryGetResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.items != null && message.hasOwnProperty("items")) {
                if (!Array.isArray(message.items))
                    return "items: array expected";
                for (let i = 0; i < message.items.length; ++i) {
                    let error = $root.ws.ItemEntry.verify(message.items[i]);
                    if (error)
                        return "items." + error;
                }
            }
            if (message.total != null && message.hasOwnProperty("total"))
                if (!$util.isInteger(message.total))
                    return "total: integer expected";
            return null;
        };

        /**
         * Creates an InventoryGetResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.InventoryGetResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.InventoryGetResp} InventoryGetResp
         */
        InventoryGetResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.InventoryGetResp)
                return object;
            let message = new $root.ws.InventoryGetResp();
            if (object.items) {
                if (!Array.isArray(object.items))
                    throw TypeError(".ws.InventoryGetResp.items: array expected");
                message.items = [];
                for (let i = 0; i < object.items.length; ++i) {
                    if (typeof object.items[i] !== "object")
                        throw TypeError(".ws.InventoryGetResp.items: object expected");
                    message.items[i] = $root.ws.ItemEntry.fromObject(object.items[i]);
                }
            }
            if (object.total != null)
                message.total = object.total | 0;
            return message;
        };

        /**
         * Creates a plain object from an InventoryGetResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.InventoryGetResp
         * @static
         * @param {ws.InventoryGetResp} message InventoryGetResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        InventoryGetResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.items = [];
            if (options.defaults)
                object.total = 0;
            if (message.items && message.items.length) {
                object.items = [];
                for (let j = 0; j < message.items.length; ++j)
                    object.items[j] = $root.ws.ItemEntry.toObject(message.items[j], options);
            }
            if (message.total != null && message.hasOwnProperty("total"))
                object.total = message.total;
            return object;
        };

        /**
         * Converts this InventoryGetResp to JSON.
         * @function toJSON
         * @memberof ws.InventoryGetResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        InventoryGetResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return InventoryGetResp;
    })();

    ws.ItemEntry = (function() {

        /**
         * Properties of an ItemEntry.
         * @memberof ws
         * @interface IItemEntry
         * @property {string|null} [itemId] ItemEntry itemId
         * @property {number|Long|null} [count] ItemEntry count
         * @property {string|null} [name] ItemEntry name
         */

        /**
         * Constructs a new ItemEntry.
         * @memberof ws
         * @classdesc Represents an ItemEntry.
         * @implements IItemEntry
         * @constructor
         * @param {ws.IItemEntry=} [properties] Properties to set
         */
        function ItemEntry(properties) {
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * ItemEntry itemId.
         * @member {string} itemId
         * @memberof ws.ItemEntry
         * @instance
         */
        ItemEntry.prototype.itemId = "";

        /**
         * ItemEntry count.
         * @member {number|Long} count
         * @memberof ws.ItemEntry
         * @instance
         */
        ItemEntry.prototype.count = $util.Long ? $util.Long.fromBits(0,0,false) : 0;

        /**
         * ItemEntry name.
         * @member {string} name
         * @memberof ws.ItemEntry
         * @instance
         */
        ItemEntry.prototype.name = "";

        /**
         * Creates a new ItemEntry instance using the specified properties.
         * @function create
         * @memberof ws.ItemEntry
         * @static
         * @param {ws.IItemEntry=} [properties] Properties to set
         * @returns {ws.ItemEntry} ItemEntry instance
         */
        ItemEntry.create = function create(properties) {
            return new ItemEntry(properties);
        };

        /**
         * Encodes the specified ItemEntry message. Does not implicitly {@link ws.ItemEntry.verify|verify} messages.
         * @function encode
         * @memberof ws.ItemEntry
         * @static
         * @param {ws.IItemEntry} message ItemEntry message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ItemEntry.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.itemId != null && Object.hasOwnProperty.call(message, "itemId"))
                writer.uint32(/* id 1, wireType 2 =*/10).string(message.itemId);
            if (message.count != null && Object.hasOwnProperty.call(message, "count"))
                writer.uint32(/* id 2, wireType 0 =*/16).int64(message.count);
            if (message.name != null && Object.hasOwnProperty.call(message, "name"))
                writer.uint32(/* id 3, wireType 2 =*/26).string(message.name);
            return writer;
        };

        /**
         * Encodes the specified ItemEntry message, length delimited. Does not implicitly {@link ws.ItemEntry.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.ItemEntry
         * @static
         * @param {ws.IItemEntry} message ItemEntry message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        ItemEntry.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes an ItemEntry message from the specified reader or buffer.
         * @function decode
         * @memberof ws.ItemEntry
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.ItemEntry} ItemEntry
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ItemEntry.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.ItemEntry();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    message.itemId = reader.string();
                    break;
                case 2:
                    message.count = reader.int64();
                    break;
                case 3:
                    message.name = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes an ItemEntry message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.ItemEntry
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.ItemEntry} ItemEntry
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        ItemEntry.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies an ItemEntry message.
         * @function verify
         * @memberof ws.ItemEntry
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        ItemEntry.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                if (!$util.isString(message.itemId))
                    return "itemId: string expected";
            if (message.count != null && message.hasOwnProperty("count"))
                if (!$util.isInteger(message.count) && !(message.count && $util.isInteger(message.count.low) && $util.isInteger(message.count.high)))
                    return "count: integer|Long expected";
            if (message.name != null && message.hasOwnProperty("name"))
                if (!$util.isString(message.name))
                    return "name: string expected";
            return null;
        };

        /**
         * Creates an ItemEntry message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.ItemEntry
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.ItemEntry} ItemEntry
         */
        ItemEntry.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.ItemEntry)
                return object;
            let message = new $root.ws.ItemEntry();
            if (object.itemId != null)
                message.itemId = String(object.itemId);
            if (object.count != null)
                if ($util.Long)
                    (message.count = $util.Long.fromValue(object.count)).unsigned = false;
                else if (typeof object.count === "string")
                    message.count = parseInt(object.count, 10);
                else if (typeof object.count === "number")
                    message.count = object.count;
                else if (typeof object.count === "object")
                    message.count = new $util.LongBits(object.count.low >>> 0, object.count.high >>> 0).toNumber();
            if (object.name != null)
                message.name = String(object.name);
            return message;
        };

        /**
         * Creates a plain object from an ItemEntry message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.ItemEntry
         * @static
         * @param {ws.ItemEntry} message ItemEntry
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        ItemEntry.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.defaults) {
                object.itemId = "";
                if ($util.Long) {
                    let long = new $util.Long(0, 0, false);
                    object.count = options.longs === String ? long.toString() : options.longs === Number ? long.toNumber() : long;
                } else
                    object.count = options.longs === String ? "0" : 0;
                object.name = "";
            }
            if (message.itemId != null && message.hasOwnProperty("itemId"))
                object.itemId = message.itemId;
            if (message.count != null && message.hasOwnProperty("count"))
                if (typeof message.count === "number")
                    object.count = options.longs === String ? String(message.count) : message.count;
                else
                    object.count = options.longs === String ? $util.Long.prototype.toString.call(message.count) : options.longs === Number ? new $util.LongBits(message.count.low >>> 0, message.count.high >>> 0).toNumber() : message.count;
            if (message.name != null && message.hasOwnProperty("name"))
                object.name = message.name;
            return object;
        };

        /**
         * Converts this ItemEntry to JSON.
         * @function toJSON
         * @memberof ws.ItemEntry
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        ItemEntry.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return ItemEntry;
    })();

    ws.BatchResp = (function() {

        /**
         * Properties of a BatchResp.
         * @memberof ws
         * @interface IBatchResp
         * @property {Array.<ws.IWSResp>|null} [msgs] BatchResp msgs
         */

        /**
         * Constructs a new BatchResp.
         * @memberof ws
         * @classdesc Represents a BatchResp.
         * @implements IBatchResp
         * @constructor
         * @param {ws.IBatchResp=} [properties] Properties to set
         */
        function BatchResp(properties) {
            this.msgs = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * BatchResp msgs.
         * @member {Array.<ws.IWSResp>} msgs
         * @memberof ws.BatchResp
         * @instance
         */
        BatchResp.prototype.msgs = $util.emptyArray;

        /**
         * Creates a new BatchResp instance using the specified properties.
         * @function create
         * @memberof ws.BatchResp
         * @static
         * @param {ws.IBatchResp=} [properties] Properties to set
         * @returns {ws.BatchResp} BatchResp instance
         */
        BatchResp.create = function create(properties) {
            return new BatchResp(properties);
        };

        /**
         * Encodes the specified BatchResp message. Does not implicitly {@link ws.BatchResp.verify|verify} messages.
         * @function encode
         * @memberof ws.BatchResp
         * @static
         * @param {ws.IBatchResp} message BatchResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        BatchResp.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.msgs != null && message.msgs.length)
                for (let i = 0; i < message.msgs.length; ++i)
                    $root.ws.WSResp.encode(message.msgs[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified BatchResp message, length delimited. Does not implicitly {@link ws.BatchResp.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.BatchResp
         * @static
         * @param {ws.IBatchResp} message BatchResp message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        BatchResp.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a BatchResp message from the specified reader or buffer.
         * @function decode
         * @memberof ws.BatchResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.BatchResp} BatchResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        BatchResp.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.BatchResp();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    if (!(message.msgs && message.msgs.length))
                        message.msgs = [];
                    message.msgs.push($root.ws.WSResp.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a BatchResp message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.BatchResp
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.BatchResp} BatchResp
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        BatchResp.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a BatchResp message.
         * @function verify
         * @memberof ws.BatchResp
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        BatchResp.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.msgs != null && message.hasOwnProperty("msgs")) {
                if (!Array.isArray(message.msgs))
                    return "msgs: array expected";
                for (let i = 0; i < message.msgs.length; ++i) {
                    let error = $root.ws.WSResp.verify(message.msgs[i]);
                    if (error)
                        return "msgs." + error;
                }
            }
            return null;
        };

        /**
         * Creates a BatchResp message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.BatchResp
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.BatchResp} BatchResp
         */
        BatchResp.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.BatchResp)
                return object;
            let message = new $root.ws.BatchResp();
            if (object.msgs) {
                if (!Array.isArray(object.msgs))
                    throw TypeError(".ws.BatchResp.msgs: array expected");
                message.msgs = [];
                for (let i = 0; i < object.msgs.length; ++i) {
                    if (typeof object.msgs[i] !== "object")
                        throw TypeError(".ws.BatchResp.msgs: object expected");
                    message.msgs[i] = $root.ws.WSResp.fromObject(object.msgs[i]);
                }
            }
            return message;
        };

        /**
         * Creates a plain object from a BatchResp message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.BatchResp
         * @static
         * @param {ws.BatchResp} message BatchResp
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        BatchResp.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.msgs = [];
            if (message.msgs && message.msgs.length) {
                object.msgs = [];
                for (let j = 0; j < message.msgs.length; ++j)
                    object.msgs[j] = $root.ws.WSResp.toObject(message.msgs[j], options);
            }
            return object;
        };

        /**
         * Converts this BatchResp to JSON.
         * @function toJSON
         * @memberof ws.BatchResp
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        BatchResp.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return BatchResp;
    })();

    ws.BatchReq = (function() {

        /**
         * Properties of a BatchReq.
         * @memberof ws
         * @interface IBatchReq
         * @property {Array.<ws.IWSReq>|null} [reqs] BatchReq reqs
         */

        /**
         * Constructs a new BatchReq.
         * @memberof ws
         * @classdesc Represents a BatchReq.
         * @implements IBatchReq
         * @constructor
         * @param {ws.IBatchReq=} [properties] Properties to set
         */
        function BatchReq(properties) {
            this.reqs = [];
            if (properties)
                for (let keys = Object.keys(properties), i = 0; i < keys.length; ++i)
                    if (properties[keys[i]] != null)
                        this[keys[i]] = properties[keys[i]];
        }

        /**
         * BatchReq reqs.
         * @member {Array.<ws.IWSReq>} reqs
         * @memberof ws.BatchReq
         * @instance
         */
        BatchReq.prototype.reqs = $util.emptyArray;

        /**
         * Creates a new BatchReq instance using the specified properties.
         * @function create
         * @memberof ws.BatchReq
         * @static
         * @param {ws.IBatchReq=} [properties] Properties to set
         * @returns {ws.BatchReq} BatchReq instance
         */
        BatchReq.create = function create(properties) {
            return new BatchReq(properties);
        };

        /**
         * Encodes the specified BatchReq message. Does not implicitly {@link ws.BatchReq.verify|verify} messages.
         * @function encode
         * @memberof ws.BatchReq
         * @static
         * @param {ws.IBatchReq} message BatchReq message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        BatchReq.encode = function encode(message, writer) {
            if (!writer)
                writer = $Writer.create();
            if (message.reqs != null && message.reqs.length)
                for (let i = 0; i < message.reqs.length; ++i)
                    $root.ws.WSReq.encode(message.reqs[i], writer.uint32(/* id 1, wireType 2 =*/10).fork()).ldelim();
            return writer;
        };

        /**
         * Encodes the specified BatchReq message, length delimited. Does not implicitly {@link ws.BatchReq.verify|verify} messages.
         * @function encodeDelimited
         * @memberof ws.BatchReq
         * @static
         * @param {ws.IBatchReq} message BatchReq message or plain object to encode
         * @param {$protobuf.Writer} [writer] Writer to encode to
         * @returns {$protobuf.Writer} Writer
         */
        BatchReq.encodeDelimited = function encodeDelimited(message, writer) {
            return this.encode(message, writer).ldelim();
        };

        /**
         * Decodes a BatchReq message from the specified reader or buffer.
         * @function decode
         * @memberof ws.BatchReq
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @param {number} [length] Message length if known beforehand
         * @returns {ws.BatchReq} BatchReq
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        BatchReq.decode = function decode(reader, length) {
            if (!(reader instanceof $Reader))
                reader = $Reader.create(reader);
            let end = length === undefined ? reader.len : reader.pos + length, message = new $root.ws.BatchReq();
            while (reader.pos < end) {
                let tag = reader.uint32();
                switch (tag >>> 3) {
                case 1:
                    if (!(message.reqs && message.reqs.length))
                        message.reqs = [];
                    message.reqs.push($root.ws.WSReq.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
                }
            }
            return message;
        };

        /**
         * Decodes a BatchReq message from the specified reader or buffer, length delimited.
         * @function decodeDelimited
         * @memberof ws.BatchReq
         * @static
         * @param {$protobuf.Reader|Uint8Array} reader Reader or buffer to decode from
         * @returns {ws.BatchReq} BatchReq
         * @throws {Error} If the payload is not a reader or valid buffer
         * @throws {$protobuf.util.ProtocolError} If required fields are missing
         */
        BatchReq.decodeDelimited = function decodeDelimited(reader) {
            if (!(reader instanceof $Reader))
                reader = new $Reader(reader);
            return this.decode(reader, reader.uint32());
        };

        /**
         * Verifies a BatchReq message.
         * @function verify
         * @memberof ws.BatchReq
         * @static
         * @param {Object.<string,*>} message Plain object to verify
         * @returns {string|null} `null` if valid, otherwise the reason why it is not
         */
        BatchReq.verify = function verify(message) {
            if (typeof message !== "object" || message === null)
                return "object expected";
            if (message.reqs != null && message.hasOwnProperty("reqs")) {
                if (!Array.isArray(message.reqs))
                    return "reqs: array expected";
                for (let i = 0; i < message.reqs.length; ++i) {
                    let error = $root.ws.WSReq.verify(message.reqs[i]);
                    if (error)
                        return "reqs." + error;
                }
            }
            return null;
        };

        /**
         * Creates a BatchReq message from a plain object. Also converts values to their respective internal types.
         * @function fromObject
         * @memberof ws.BatchReq
         * @static
         * @param {Object.<string,*>} object Plain object
         * @returns {ws.BatchReq} BatchReq
         */
        BatchReq.fromObject = function fromObject(object) {
            if (object instanceof $root.ws.BatchReq)
                return object;
            let message = new $root.ws.BatchReq();
            if (object.reqs) {
                if (!Array.isArray(object.reqs))
                    throw TypeError(".ws.BatchReq.reqs: array expected");
                message.reqs = [];
                for (let i = 0; i < object.reqs.length; ++i) {
                    if (typeof object.reqs[i] !== "object")
                        throw TypeError(".ws.BatchReq.reqs: object expected");
                    message.reqs[i] = $root.ws.WSReq.fromObject(object.reqs[i]);
                }
            }
            return message;
        };

        /**
         * Creates a plain object from a BatchReq message. Also converts values to other types if specified.
         * @function toObject
         * @memberof ws.BatchReq
         * @static
         * @param {ws.BatchReq} message BatchReq
         * @param {$protobuf.IConversionOptions} [options] Conversion options
         * @returns {Object.<string,*>} Plain object
         */
        BatchReq.toObject = function toObject(message, options) {
            if (!options)
                options = {};
            let object = {};
            if (options.arrays || options.defaults)
                object.reqs = [];
            if (message.reqs && message.reqs.length) {
                object.reqs = [];
                for (let j = 0; j < message.reqs.length; ++j)
                    object.reqs[j] = $root.ws.WSReq.toObject(message.reqs[j], options);
            }
            return object;
        };

        /**
         * Converts this BatchReq to JSON.
         * @function toJSON
         * @memberof ws.BatchReq
         * @instance
         * @returns {Object.<string,*>} JSON object
         */
        BatchReq.prototype.toJSON = function toJSON() {
            return this.constructor.toObject(this, $protobuf.util.toJSONOptions);
        };

        return BatchReq;
    })();

    return ws;
})();

export { $root as default };
