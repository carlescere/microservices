####
# Assumptions:
# -> Image is created after building process has taken place
# -> The bin folder exists as product of said building process.
####

FROM progrium/busybox

MAINTAINER Carles Cerezo <carlescere@gmail.com>

COPY bin/movies /movies
COPY database /database

EXPOSE 8001

CMD ["/movies","-port", "8001"]
