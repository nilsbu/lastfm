# Set the base image to Node.js 18 (alpine)
FROM node:18-alpine

# Set the working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the current directory contents to the container at /app
COPY . .

RUN npm run build

# Expose port 3000 for the app to listen on
EXPOSE 3000

# Start the app with npm start
CMD ["npm", "start"]
