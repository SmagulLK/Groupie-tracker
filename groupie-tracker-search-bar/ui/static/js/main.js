var card = document.querySelector(".card")
var relations = card.querySelector(".relations");
var members = card.querySelector(".members")
var moreRelationsBtn = card.querySelector("#more");
var membersBtn = card.querySelector("#member");
var closeRelationsBtn = card.querySelector("#close");
var closeMembersBtn = card.querySelector("#close-members");


let location2 = document.getElementsByClassName('loc');
let arr = []

for (key of location2) {
    arr.push(key.innerHTML.replace(/[^\w\s]|_/g," "));
}

console.log(location2)

ymaps.ready(init);

function init () {
   myMap = new ymaps.Map('map', {  
        center: [30, 10],
        zoom: 1
    });

    for (key of arr) {
        var myGeocoder = ymaps.geocode(key, {results: 1, prefLang: "en"});
        console.log(key)
        myGeocoder.then(
            function (res) {
                myMap.geoObjects.add(res.geoObjects);  
            },
        );
    }

}

moreRelationsBtn.addEventListener("click", () => {
    relations.classList.toggle("show");
});
closeRelationsBtn.addEventListener("click", () => {
    moreRelationsBtn.click();
})

membersBtn.addEventListener("click", () => {
    members.classList.toggle("show-members");
});

closeMembersBtn.addEventListener("click", () => {
    membersBtn.click();
})

const accordion = document.getElementsByClassName('contentBx')
    for (i = 0; i < accordion.length; i++) {
        accordion[i].addEventListener('click', function(){
             this.classList.toggle('active')              
            })            
        } 


