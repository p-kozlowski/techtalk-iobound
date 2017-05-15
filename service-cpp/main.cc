/*
 * This file is open source software, licensed to you under the terms
 * of the Apache License, Version 2.0 (the "License").  See the NOTICE file
 * distributed with this work for additional information regarding copyright
 * ownership.  You may not use this file except in compliance with the License.
 *
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
/*
 * Copyright 2015 Cloudius Systems
 */

#include <functional>
#include "http/httpd.hh"
#include "http/handlers.hh"
#include "http/function_handlers.hh"
#include "http/file_handler.hh"
#include "apps/httpd/demo.json.hh"
#include "http/api_docs.hh"

namespace bpo = boost::program_options;

using namespace httpd;

class UserHandler : public httpd::handler_base {
public:
    UserHandler(const std::string &addr, const uint16_t port) : dataStoreAddr(addr, port) {

    }

    future<std::unique_ptr<reply>> handle(const sstring& path,
            std::unique_ptr<request> req, std::unique_ptr<reply> rep) override {
        auto addr = ipv4_addr("127.0.0.1", 3333);
        return connect(addr).then([req = std::move(req), rep = std::move(rep)](connected_socket s) mutable {
            auto name = req->param["name"];
            auto out = s.output();
            return do_with(std::move(s), std::move(out), std::move(name), [rep = std::move(rep)](auto &s, auto &out, auto &name) mutable {
                return out.write(name+sstring("\n")).then([&s, &out, rep = std::move(rep)]() mutable {
                    return out.close().then([&s, rep = std::move(rep)]() mutable {
                        return s.input().read().then([rep = std::move(rep)](auto data) mutable {
                            rep->_content = sstring("Hello, ") + sstring(data.get(), data.size());
                            rep->done("html");
                            return make_ready_future<std::unique_ptr<reply>>(std::move(rep));
                        });
                    });
                });
            });
        });
    }
private:
    const ipv4_addr dataStoreAddr;
};


int main(int ac, char** av) {
    app_template app;
    app.add_options()("port", bpo::value<uint16_t>()->default_value(10000), "HTTP Server port");
    app.add_options()("datastore-port", bpo::value<uint16_t>()->default_value(3333), "Datastore port");
    app.add_options()("datastore-addr", bpo::value<std::string>()->default_value("127.0.0.1"), "Datastore address");

    return app.run_deprecated(ac, av, [&] {
        auto&& config = app.configuration();
        const uint16_t port = config["port"].as<uint16_t>();

        auto server = new http_server_control();
        server->start().then([server, &config] {
            const uint16_t datastorePort = config["datastorePort"].as<uint16_t>();
            const std::string datastoreAddr = config["datastoreAddr"].as<std::string>();
            return server->set_routes([&datastoreAddr, &datastorePort](routes& r) {
                r.add(operation_type::GET, url("/users").remainder("name"), new UserHandler{datastoreAddr, datastorePort});
            });
        }).then([server, port] {
            return server->listen(port);
        }).then([server, port] {
            std::cout << "Seastar HTTP server listening on port " << port << " ...\n";
            engine().at_exit([server] {
                return server->stop();
            });
        });
    });
}
