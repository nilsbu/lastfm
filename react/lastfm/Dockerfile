# Set the base image to Node.js 18 (alpine)
FROM node:18-alpine

ENV BACKEND=lastfm-srv:3000

# Set the working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json to the container
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the app's source code to the container
COPY . .

# Expose port 3000 for the app to listen on
EXPOSE 3000

# Start the app with npm start
CMD ["npm", "start"]
