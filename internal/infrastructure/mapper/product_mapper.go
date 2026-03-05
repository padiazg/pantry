package mapper

// TODO: import domain entities when mapping is implemented
// "github.com/padiazg/pantry/internal/core/domain"

// ProductMapper Mapper para convertir entre la entidad Product del dominio y DTOs/modelos de base de datos
type ProductMapper struct {
	// TODO: Add dependencies if needed
}

// NewProductMapper creates a new ProductMapper instance
func NewProductMapper() *ProductMapper {
	return &ProductMapper{}
}

// ToDTO converts a domain entity to a DTO
func (m *ProductMapper) ToDTO(entity interface{}) interface{} {
	// TODO: Implement domain to DTO mapping
	// Example:
	// if user, ok := entity.(*domain.User); ok {
	//     return &UserDTO{
	//         ID:    user.ID,
	//         Name:  user.Name,
	//         Email: user.Email,
	//     }
	// }
	return nil
}

// ToDomain converts a DTO to a domain entity
func (m *ProductMapper) ToDomain(dto interface{}) (interface{}, error) {
	// TODO: Implement DTO to domain mapping
	// Example:
	// if userDTO, ok := dto.(*UserDTO); ok {
	//     return domain.NewUser(
	//         userDTO.ID,
	//         userDTO.Name,
	//         userDTO.Email,
	//     )
	// }
	return nil, nil
}

// ToListDTO converts a slice of domain entities to DTOs
func (m *ProductMapper) ToListDTO(entities []interface{}) []interface{} {
	dtos := make([]interface{}, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, m.ToDTO(entity))
	}
	return dtos
}
