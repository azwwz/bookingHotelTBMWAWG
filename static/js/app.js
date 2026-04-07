// Prompt is our JavaScript module for all alerts, notifications, and custom popup dialog
function Prompt() {
    let toast = function (c) {
        const {
            msg = '',
            icon = 'success',
            position = 'top-end',

        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })

    }

    // define a func for response availability
    async function custom(c) {
        const {
            msg = "",
            title = "",
            showConfirmButton = true,
        } = c;

        // construct detail for availability
        const { value: result } = await Swal.fire({
            title: title,
            html: msg,
            backdrop: false,
            allowOutsideClick:false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton : showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined){
                    c.willOpen()
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined){
                    c.didOpen()
                }
            },
            preConfirm: () => {
                if (c.preConfirm !== undefined){
                    c.preConfirm()
                }
            }
        })

        // deal result after dialog end
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== ""){
                    if (c.callback !== undefined){
                        // use c's callback
                        c.callback(result)
                    }
                } else{
                    c.callback(false)
                }
            }else{
                c.callback(false)
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

function generalAndMajors(rc){
    document.getElementById("check-availability-button").addEventListener("click", function () {
        let html = `
        <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
            <div class="form-row">
                <div class="col">
                    <div class="form-row" id="reservation-dates-modal">
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                        </div>
                        <div class="col">
                            <input disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                        </div>

                    </div>
                </div>
            </div>
        </form>
        `;
        attention.custom({
            title: 'Choose your dates',
            msg: html,
            showConfirmButton: true,
            callback: function(result){
                console.log("called");
                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token",rc.csrf)
                formData.append("room_id",rc.r_id)
                fetch("/search-availability-json",{
                    method:"post",
                    body:formData,
                })
                    .then(Response => Response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: 'success',
                                showConfirmButton: false,

                                msg: `<p>Room is available!</p>
                                  <p><a href="/book-room?id=${data.room_id}&s=${data.start_date}&e=${data.end_date}"
                                      class="btn btn-primary">
                                      Book Now!
                                  </a></p>`,
                            })
                            console.log(data.room_id)
                        }else{
                            attention.error({
                                msg : "no availability"
                            })
                        }
                    })
            },
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: new Date(),
                })
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value
                ]
            }
        });
    })

}