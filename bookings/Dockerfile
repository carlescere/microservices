####
# Assumptions:
# -> Image is created after building process has taken place
# -> The bin folder exists as product of said building process.
####

FROM progrium/busybox

MAINTAINER Carles Cerezo <carlescere@gmail.com>

COPY bin/bookings /bookings
COPY database /database

EXPOSE 8003

CMD ["/bookings","-port", "8003"]
