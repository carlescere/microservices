####
# Assumptions:
# -> Image is created after building process has taken place
# -> The bin folder exists as product of said building process.
####

FROM progrium/busybox

MAINTAINER Carles Cerezo <carlescere@gmail.com>

COPY bin/users /users
COPY database /database

EXPOSE 8000

CMD ["/users","-port", "8001"]
