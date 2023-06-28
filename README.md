# Poosible Backend 🌟

Welcome to the Poosible Backend repository! This repository contains the backend code for the Poosible project. Poosible is a Point of Sale (POS) system designed to cater to the needs of small businesses. It offers inventory management and transaction organization features to empower small business owners in handling their operations effectively.

## 💡 Key Features

Inventory Management: With Poosible, businesses can easily manage their inventory by keeping track of available products, monitoring stock levels, and receiving alerts for low stock or out-of-stock items. This helps businesses stay informed and ensures they can meet customer demands effectively.

Item Organization: Poosible allows businesses to organize their items efficiently. It enables the creation of categories, tags, and attributes for items, making it easier to navigate and search for specific products. This feature enhances the overall organization and accessibility of products within the system.

Customer Management: Poosible enables businesses to maintain a database of customers, including their contact information and purchase history. This feature helps businesses build stronger relationships with customers by providing personalized experiences and targeted marketing campaigns.

Reporting and Analytics: Poosible offers comprehensive reporting and analytics tools that provide valuable insights into business performance. It generates reports on sales, inventory, and customer data, enabling business owners to make data-driven decisions and identify areas for improvement.

## 🚀 Getting Started

To get started with the Poosible Backend, follow these steps:

1. Clone this repository: `git clone https://github.com/EmirWorld/poosible-backend.git`
2. Install the required dependencies: `go mod init`
4. Start the server: `go run main.go`

That's it! You now have the Poosible Backend up and running locally on your machine.

## 📁 Project Structure

```
├── controllers/      # Contains the controllers for different API routes
├── models/           # Defines the database models using Mongoose
├── routes/           # Defines the API routes and their corresponding controllers
├── config/           # Contains configuration for database and swagger
├── utils/            # Utility functions and helpers
├── main.go           # Entry point of the application
└── ...
```

## 🛠️ Technologies Used

- Golang
- Gin Gonic Framework
- MongoDB
- Mongoose

## 📝 API Documentation

The Poosible Backend provides the following API endpoints:
- `/api/auth` - Auth-related endpoints (login, logout)
- `/api/users` - User-related endpoints (register, update profile, etc.)
- `/api/projects` - Project-related endpoints (create, update, delete, etc.)
- `/api/tasks` - Task-related endpoints (create, update, delete, etc.)
- `/api/comments` - Comment-related endpoints (create, update, delete, etc.)

For detailed documentation on each endpoint, please refer to the [API Documentation](API_DOCUMENTATION.md).

## 🐳 Docker

To build and run the Poosible Backend using Docker, we have provided a `Makefile` with the following commands:

- `make build`: Builds the Docker containers using `docker-compose build`.
- `make start`: Starts the Docker containers in detached mode using `docker-compose up -d`.

To use these commands, make sure you have Docker and Docker Compose installed on your machine. Then, simply run `make build` to build the containers and `make start` to start the containers.

## 💡 Contributing

Contributions are always welcome! If you'd like to contribute to the Poosible Backend, please follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b my-feature-branch`.
3. Make your changes and commit them: `git commit -m 'Add some feature'`.
4. Push to the branch: `git push origin my-feature-branch`.
5. Submit a pull request.

## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 📧 Contact

If you have any questions or suggestions regarding the Poosible Backend, feel free to reach out to us at [Emir Kovacevic](mailto:emirkovacevic@protonmail.com).

Let's make amazing things Poosible together! ✨🙌
