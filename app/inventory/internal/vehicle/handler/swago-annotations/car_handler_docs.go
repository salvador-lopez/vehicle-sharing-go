package swago_annotations

// GetCar godoc
// @Summary      Get a car by ID
// @Description  Returns a car resource with decoded VIN data
// @Tags         car
// @Produce      json
// @Param        id   path      string  true  "Car UUID"
// @Success      200  {object}  projection.Car
// @Failure      400  {object}  rest.ErrorResponse
// @Failure      404  {object}  rest.ErrorResponse
// @Failure      500  {object}  rest.ErrorResponse
// @Router       /cars/{id} [get]
func GetCar() {}

// CreateCar godoc
// @Summary      Create a new car
// @Description  Creates a new car record
// @Tags         car
// @Accept       json
// @Produce      json
// @Param        car  body      command.CreateCar true  "Create Car Body Params"
// @Success      201  {string}  string            "Created"
// @Failure      400  {object}  rest.ErrorResponse
// @Failure      409  {object}  rest.ErrorResponse
// @Failure      500  {object}  rest.ErrorResponse
// @Router       /cars [post]
func CreateCar() {}
