package models

import "time"

type Barrel struct {
	Id               int     `json:"id"`
	Name             string    `json:"name"`
	VolumeLiters     float64   `json:"volume_liters"`
	CurrentVolume    float64   `json:"current_volume"`
	LocationName     string    `json:"location_name"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	AssignedSellerId *int    `json:"assigned_seller_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}


type CreateBarrel struct {
	Name             string    `json:"name"`
	VolumeLiters     float64   `json:"volume_liters"`
	CurrentVolume    float64   `json:"current_volume"`
	LocationName     string    `json:"location_name"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	AssignedSellerId *int    `json:"assigned_seller_id"`
}

type UpdateBarrel struct {
	Id               int     `json:"id"`
	Name             string    `json:"name"`
	VolumeLiters     float64   `json:"volume_liters"`
	CurrentVolume    float64   `json:"current_volume"`
	LocationName     string    `json:"location_name"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	AssignedSellerId *int    `json:"assigned_seller_id"`
}

type GetListBarrelRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListBarrelResponse struct {
	Count int     `json:"count"`
	Barrels []*Barrel `json:"Barrels"`
}
