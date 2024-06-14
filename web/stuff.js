function getCyclists() {
    fetch('http://localhost:8080/getCyclists')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            appendCyclists(data);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
}

function appendCyclists(cyclists) {
    const databaseDiv = document.getElementById('database');
    databaseDiv.innerHTML = ''; // Clear existing content
    const title = document.createElement('h2');
    title.innerText = 'Table: Cyclists';
    databaseDiv.appendChild(title);
    databaseDiv.appendChild(document.createElement('hr'));
    cyclists.forEach(cyclist => {
        const cyclistInfo = document.createElement('div');
        cyclistInfo.innerHTML = `
            <p>ID: ${cyclist.id}</p>
            <p>First Name: ${cyclist.name}</p>
            <p>Phone Number: ${cyclist.phone_number}</p>
            <p>Skill Level: ${cyclist.skill_level}</p>
            <p>Address ID: ${cyclist.address_id}</p>
            <p>Bike ID: ${cyclist.bike_id}</p>
            <hr>
        `;
        databaseDiv.appendChild(cyclistInfo);
    });
}

function getBikes() {
    fetch('http://localhost:8080/getBikes')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            appendBikes(data);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
}

function appendBikes(bikes) {
    const databaseDiv = document.getElementById('database');
    const title = document.createElement('h2');
    title.innerText = 'Table: Bikes';
    databaseDiv.appendChild(title);
    databaseDiv.appendChild(document.createElement('hr'));
    bikes.forEach(bikes => {
        const bikesInfo = document.createElement('div');
        bikesInfo.innerHTML = `
            <p>ID: ${bikes.id}</p>
            <p>Nickname: ${bikes.nickname}</p>
            <p>Serial Number: ${bikes.serial_number}</p>
            <p>Year: ${bikes.year}</p>
            <p>Model: ${bikes.model}</p>
            <p>Make: ${bikes.make}</p>
            <p>Mileage: ${bikes.mileage}</p>
            <hr>
        `;
        databaseDiv.appendChild(bikesInfo);
    });
}

function getAddresses() {
    fetch('http://localhost:8080/getAddresses')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.json();
        })
        .then(data => {
            appendAddresses(data);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
}

function appendAddresses(addresses) {
    const databaseDiv = document.getElementById('database');
    const title = document.createElement('h2');
    title.innerText = 'Table: Addresses';
    databaseDiv.appendChild(title);
    databaseDiv.appendChild(document.createElement('hr'));
    addresses.forEach(addresses => {
        const addressesInfo = document.createElement('div');
        addressesInfo.innerHTML = `
            <p>ID: ${addresses.id}</p>
            <p>Street: ${addresses.street}</p>
            <p>State: ${addresses.state}</p>
            <p>Zip Code: ${addresses.zip}</p>
            <hr>
        `;
        databaseDiv.appendChild(addressesInfo);
    });
}


async function submitQuestion() {
    const inputElement = document.getElementById("questionInput");
    const question = inputElement.value;
    
    // Make the input field blank
    inputElement.value = "";
  
    const url = "http://localhost:8080/ask"; // Replace with your server's URL if different
  
    const requestBody = {
      text: question
    };
  
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(requestBody)
      });
  
      if (!response.ok) {
        throw new Error("Network response was not ok " + response.statusText);
      }
  
      const data = await response.json();
      console.log("Response from server:", data.response);
  
      // Display the response in the responseContainer div
      const outerContainer = document.getElementById("outputContainer");
      //make its display visible
        outerContainer.style.visibility = "visible";
      const responseContainer = document.getElementById("output");
      responseContainer.innerText = data.response;
      setTimeout(() => {
        getCyclists();
        setTimeout(() => {
            getBikes();
            getAddresses();
            } , 250);
        }, 1000);
    } catch (error) {
      console.error("There was a problem with the fetch operation:", error);
    }
  }

  document.addEventListener("DOMContentLoaded", function() {
    // Call the function to fetch and display cyclists
    getCyclists();
    setTimeout(() => {
        getBikes();
        getAddresses();
        } , 250);
});