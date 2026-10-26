export type EntityID = string

export interface FabricColor {
  name: string
  hex: string
}

export interface Fabric {
  id: EntityID
  ownerId: EntityID
  name: string
  material: string
  pattern: string
  styleTags: string[]
  colors: FabricColor[]
  coverAssetId?: EntityID
  memo?: string
  version: number
  createdAt: string
  updatedAt: string
}

export interface FabricPiece {
  id: EntityID
  fabricId: EntityID
  label: string
  lengthMm?: number
  widthMm?: number
  storageLocation?: string
  status: 'in_stock' | 'used_up' | 'archived'
  version: number
}

export interface Project {
  id: EntityID
  ownerId: EntityID
  name: string
  kind: string
  status: 'planned' | 'making' | 'completed' | 'archived'
  memo?: string
  version: number
  createdAt: string
  updatedAt: string
}

