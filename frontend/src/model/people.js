/*

Copyright (c) 2018 - 2021 Michael Mayer <hello@photoprism.org>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

    PhotoPrism® is a registered trademark of Michael Mayer.  You may use it as required
    to describe our software, run your own server, for educational purposes, but not for
    offering commercial goods, products, or services without prior written permission.
    In other words, please ask.

Feel free to send an e-mail to hello@photoprism.org if you have questions,
want to support our work, or just want to say hello.

Additional information can be found in our Developer Guide:
https://docs.photoprism.org/developer-guide/

*/

import RestModel from "model/rest";
import Api from "common/api";
import {DateTime} from "luxon";
import {config} from "../session";
import {$gettext} from "common/vm";

export class People extends RestModel {
    getDefaults() {
        return {
            UID: "",
            FullName: "",
            UserId: null,
            BoD: null,
            DeadDate: null,
            PhotoCount: 0,
            PlaceCount: 0,
            CreatedAt: "",
            UpdatedAt: "",
            DelateAt: "",
        };
    }

    getEntityName() {
        return this.Slug;
    }

    getFullName() {
        return this.FullName;
    }

    /*thumbnailUrl(size) {
        return `/api/v1/peoples/${this.getId()}/t/${config.previewToken()}/${size}`;
    }*/

    getDoB() {
        let date = this.BoD

        return DateTime.fromISO(`${date}T12:00:00Z`).toUTC();
    }

    localDoB(time) {
        if (!this.TakenAtLocal) {
            return this.utcDate();
        }

        let zone = this.getTimeZone();

        return DateTime.fromISO(this.localDateString(time), {zone});
    }

    getDoBString() {
        return this.localDate().toLocaleString(DateTime.DATE_HUGE);
    }

    getCreatedString() {
        return DateTime.fromISO(this.CreatedAt).toLocaleString(DateTime.DATETIME_MED);
    }


    static getCollectionResource() {
        return "peoples";
    }

    static getModelName() {
        return $gettext("Peoples");
    }
}

export default People;
