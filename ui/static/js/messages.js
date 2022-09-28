"use strict";

document.addEventListener("DOMContentLoaded", () => {
  createTable();
});

async function createTable() {
  try {
    const url = "/json";
    const result = await fetch(url, {
      headers: { "Content-type": "application/json" },
    });
    const db = await result.json();
    let i = 1;
    db.forEach((element) => {
      const { Name, Email, Title, Message } = element;
      const nameEl = document.createElement("TD");
      nameEl.textContent = Name;

      const emailEl = document.createElement("TD");
      emailEl.textContent = Email;

      const titleEl = document.createElement("TD");
      titleEl.textContent = Title;

      const msgEl = document.createElement("TD");
      msgEl.textContent = Message;

      const rowEl = document.createElement("TH");
      rowEl.textContent = i;
      rowEl.setAttribute("scope", "row");
      i++;

      const rowContainer = document.createElement("TR");
      rowContainer.appendChild(rowEl);
      rowContainer.appendChild(nameEl);
      rowContainer.appendChild(emailEl);
      rowContainer.appendChild(titleEl);
      rowContainer.appendChild(msgEl);

      document.querySelector("tbody").appendChild(rowContainer);
    });
    console.log("workin");
  } catch (error) {
    console.log(error);
  }
}
