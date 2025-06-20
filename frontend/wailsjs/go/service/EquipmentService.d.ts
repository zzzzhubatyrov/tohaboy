// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {model} from '../models';

export function CreateEquipment(arg1:model.Equipment):Promise<model.EquipmentResponse>;

export function DeleteEquipment(arg1:number):Promise<model.EquipmentResponse>;

export function GetAllEquipment():Promise<model.EquipmentListResponse>;

export function GetEquipment(arg1:number):Promise<model.EquipmentResponse>;

export function GetEquipmentByLocation(arg1:number):Promise<model.EquipmentListResponse>;

export function GetEquipmentBySupplier(arg1:number):Promise<model.EquipmentListResponse>;

export function UpdateEquipment(arg1:model.Equipment):Promise<model.EquipmentResponse>;
